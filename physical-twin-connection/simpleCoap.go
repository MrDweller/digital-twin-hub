package physicaltwinconnection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MrDweller/digital-twin-hub/models"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/udp"
)

type SimpleCoapConnectionModel struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (simpleCoapConnectionModel SimpleCoapConnectionModel) connect() error {
	return nil
}

func (simpleCoapConnectionModel SimpleCoapConnectionModel) HandleControllCommand(controllCommandModel models.ControllCommandModel, commands any) (map[string]any, error) {
	target := fmt.Sprintf("%s:%d", simpleCoapConnectionModel.Address, simpleCoapConnectionModel.Port)
	co, err := udp.Dial(target)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// requestPayload := fmt.Sprintf("%v", commands)

	requestPayload, err := json.Marshal(commands)
	if err != nil {
		return nil, err
	}

	resp, err := co.Post(ctx, controllCommandModel.ServiceUri, message.TextPlain, bytes.NewReader(requestPayload))
	if err != nil {
		return nil, err
	}
	log.Printf("Response: %v", resp.String())
	payload, err := resp.ReadBody()
	if err != nil {
		return nil, err
	}
	log.Printf("Response payload: %v", string(payload))

	var response map[string]any
	err = json.Unmarshal(payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (simpleCoapConnectionModel SimpleCoapConnectionModel) HandleSensorRequest(sensedPropertyModel models.SensedPropertyModel) (map[string]any, error) {
	target := fmt.Sprintf("%s:%d", simpleCoapConnectionModel.Address, simpleCoapConnectionModel.Port)
	co, err := udp.Dial(target)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := co.Get(ctx, sensedPropertyModel.ServiceUri)
	if err != nil {
		return nil, err
	}
	log.Printf("Response: %v", resp.String())
	payload, err := resp.ReadBody()
	if err != nil {
		return nil, err
	}
	log.Printf("Response payload: %v", string(payload))

	var response map[string]any
	err = json.Unmarshal(payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
