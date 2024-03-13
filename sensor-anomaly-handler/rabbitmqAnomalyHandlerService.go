package sensoranomalyhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const EXCHANGE = "exchange"

type RabbitmqAnomalyHandlerService struct {
	rabbitmqAddress string
	rabbitmqPort    int
}

type RabbitMQHandleableAnomaly struct {
	HandleableAnomalyBase
	id uuid.UUID
}

func (rabbitMQHandleableAnomaly *RabbitMQHandleableAnomaly) GetMetaData() map[string]string {
	return map[string]string{
		EXCHANGE: fmt.Sprintf("%s-%s", rabbitMQHandleableAnomaly.EventType, rabbitMQHandleableAnomaly.id.String()),
	}
}

func (rabbitMQHandleableAnomaly *RabbitMQHandleableAnomaly) GetAnomaly() Anomaly {
	return rabbitMQHandleableAnomaly.Anomaly
}

func (service RabbitmqAnomalyHandlerService) HandleAnomaly(handleableAnomaly HandleableAnomaly, work Work) error {
	RabbitMQHandleableAnomaly := handleableAnomaly.(*RabbitMQHandleableAnomaly)
	metadata := RabbitMQHandleableAnomaly.GetMetaData()
	err := service.emit(handleableAnomaly.GetAnomaly(), metadata[EXCHANGE], work)
	return err
}

func (service RabbitmqAnomalyHandlerService) emit(anomaly Anomaly, exchange string, work Work) error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%d/", service.rabbitmqAddress, service.rabbitmqPort))
	if err != nil {
		log.Printf("%s: %s", "Failed to connect to RabbitMQ", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%s: %s", "Failed to open a channel", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Printf("%s: %s", "Failed to declare an exchange", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventDTO := EventDTO{
		Anomaly: anomaly,
		Work:    work,
	}
	data, err := json.Marshal(eventDTO)
	if err != nil {
		log.Printf("%s: %s", "Failed to marshal the anomaly", err)
		return err
	}
	err = ch.PublishWithContext(ctx,
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)

	if err != nil {
		log.Printf("%s: %s", "Failed to publish a message", err)
		return err
	}

	log.Printf(" [x] Sent %s", anomaly.EventType)
	return nil
}
