package sensoranomalyhandler

import (
	"fmt"

	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	"github.com/MrDweller/service-registry-connection/models"
)

type HandleableAnomaly struct {
	anomalyHandlingSystem models.SystemDefinition
	Anomaly
}

func (anomaly *HandleableAnomaly) GetService() additionalservice.AdditionalServiceModel {
	return additionalservice.AdditionalServiceModel{
		ServiceDefinition: models.ServiceDefinition{
			ServiceDefinition: anomaly.AnomalyType,
			ServiceUri:        fmt.Sprintf("/handle/%s", anomaly.AnomalyType),
		},
		Interfaces: []string{
			"AMQP-INSECURE-JSON",
		},
		HasExternalHost:  true,
		SystemDefinition: anomaly.anomalyHandlingSystem,
	}
}

func (anomaly *HandleableAnomaly) setAnomalyHandlingSystem(externalSystem models.SystemDefinition) {
	anomaly.anomalyHandlingSystem = externalSystem
}
