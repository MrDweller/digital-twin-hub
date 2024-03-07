package additionalservice

import (
	serviceModels "github.com/MrDweller/service-registry-connection/models"
)

type AdditionalServiceModel struct {
	ServiceDefinition serviceModels.ServiceDefinition
	Interfaces        []string
	Metadata          map[string]string
	HasExternalHost   bool
	SystemDefinition  serviceModels.SystemDefinition
}

type AdditionalService interface {
	GetService() AdditionalServiceModel
	GetMetaData() map[string]string
}
