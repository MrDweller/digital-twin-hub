package sensoranomalyhandler

import (
	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	"github.com/MrDweller/service-registry-connection/models"
)

type HandleableAnomalyBase struct {
	anomalyHandlingSystem models.SystemDefinition
	Anomaly
}

func (handleableAnomalyBase *HandleableAnomalyBase) GetService() additionalservice.AdditionalServiceModel {
	return additionalservice.AdditionalServiceModel{
		ServiceDefinition: models.ServiceDefinition{
			ServiceDefinition: handleableAnomalyBase.EventType,
			ServiceUri:        "",
		},
		Interfaces: []string{
			"AMQP-INSECURE-JSON",
		},
		HasExternalHost:  true,
		SystemDefinition: handleableAnomalyBase.anomalyHandlingSystem,
	}
}

func (handleableAnomalyBase *HandleableAnomalyBase) GetMetaData() map[string]string {
	return map[string]string{}
}

func (handleableAnomalyBase *HandleableAnomalyBase) GetAnomaly() Anomaly {
	return handleableAnomalyBase.Anomaly
}

type HandleableAnomaly interface {
	GetMetaData() map[string]string
	GetAnomaly() Anomaly
}
