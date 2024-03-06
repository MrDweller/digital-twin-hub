package additionalservice

import (
	serviceModels "github.com/MrDweller/service-registry-connection/models"
)

type AdditionalServiceModel struct {
	ServiceDefinition serviceModels.ServiceDefinition
	Interfaces        []string
	HasExternalHost   bool
	SystemDefinition  serviceModels.SystemDefinition
}

type AdditionalService interface {
	GetService() AdditionalServiceModel
}
