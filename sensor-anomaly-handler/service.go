package sensoranomalyhandler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqAnomalyHandlerService struct {
}

func (service RabbitmqAnomalyHandlerService) HandleAnomaly(anomaly Anomaly) error {
	err := service.Emit(anomaly)
	return err
}

func (service RabbitmqAnomalyHandlerService) Emit(anomaly Anomaly) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
		"logs",   // name
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

	data, err := json.Marshal(anomaly)
	if err != nil {
		log.Printf("%s: %s", "Failed to marshal the anomaly", err)
		return err
	}
	err = ch.PublishWithContext(ctx,
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)

	if err != nil {
		log.Printf("%s: %s", "Failed to publish a message", err)
		return err
	}

	log.Printf(" [x] Sent %s", anomaly.AnomalyType)
	return nil
}