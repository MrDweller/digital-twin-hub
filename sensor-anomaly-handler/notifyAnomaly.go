package sensoranomalyhandler

import (
	"fmt"

	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	"github.com/MrDweller/service-registry-connection/models"
)

type NotifyAnomaly struct {
	Anomaly
}

func (notifyAnomaly *NotifyAnomaly) GetService() additionalservice.AdditionalServiceModel {
	return additionalservice.AdditionalServiceModel{
		ServiceDefinition: models.ServiceDefinition{
			ServiceDefinition: notifyAnomaly.AnomalyType,
			ServiceUri:        fmt.Sprintf("/notify/%s", notifyAnomaly.AnomalyType),
		},
		Interfaces: []string{
			"HTTP-SECURE-JSON",
		},
		HasExternalHost: false,
	}
}

func (notifyAnomaly *NotifyAnomaly) GetMetaData() map[string]string {
	return map[string]string{}
}
