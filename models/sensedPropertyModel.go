package models

import (
	"github.com/MrDweller/service-registry-connection/models"
)

type SensedPropertyModel struct {
	models.ServiceDefinition
	SensorEndpointMode SensorEndpointMode `json:"sensorEndpointMode"`
	IntervalTime       int                `json:"intervalTime"`
}

type SensorEndpointMode string

const (
	IMMEDIATE_RETRIEVAL = "IMMEDIATE_RETRIEVAL"
	INTERVAL_RETRIEVAL  = "INTERVAL_RETRIEVAL"
)
