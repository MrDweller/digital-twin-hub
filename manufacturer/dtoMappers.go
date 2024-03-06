package manufacturer

import (
	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	"github.com/MrDweller/digital-twin-hub/models"
	sensoranomalyhandler "github.com/MrDweller/digital-twin-hub/sensor-anomaly-handler"
	serviceModels "github.com/MrDweller/service-registry-connection/models"
)

func mapDigitalTwinDtoToDigitalTwinModel(digitalTwinDTO DigitalTwinDTO, additionalServices []additionalservice.AdditionalService) models.DigitalTwinModel {
	additionalServiceModels := []additionalservice.AdditionalServiceModel{}
	for _, additionalService := range additionalServices {
		additionalServiceModel := additionalService.GetService()
		additionalServiceModels = append(additionalServiceModels, additionalServiceModel)

	}

	return models.DigitalTwinModel{
		SensedProperties:        mapSensedPropertiesDtoToSensedPropertiesModel(digitalTwinDTO.SensedProperties),
		ControlCommands:         mapControllCommandsDtoToControllCommandsModel(digitalTwinDTO.ControlCommands),
		AdditionalServiceModels: additionalServiceModels,

		PhysicalTwinConnectionModel: mapConnectionDtoToPhysicalTwinConnectionModel(digitalTwinDTO.PhysicalTwinConnection),
	}
}

func mapConnectionDtoToPhysicalTwinConnectionModel(connectionDTO ConnectionDTO) models.PhysicalTwinConnectionModel {
	return models.PhysicalTwinConnectionModel{
		ConnectionType:  models.PhysicalTwinConnectionType(connectionDTO.ConnectionType),
		ConnectionModel: connectionDTO.ConnectionModel,
	}
}

func mapHandleableAnomaliesDtoToHandleableAnomalies(handleableAnomaliesDTO []HandleableAnomalyDTO) []sensoranomalyhandler.HandleableAnomaly {
	handleableAnomalies := []sensoranomalyhandler.HandleableAnomaly{}
	for _, handleableAnomalyDTO := range handleableAnomaliesDTO {
		handleableAnomalies = append(handleableAnomalies, mapHandleableAnomalyDtoToHandleableAnomaly(handleableAnomalyDTO))
	}
	return handleableAnomalies
}

func mapHandleableAnomalyDtoToHandleableAnomaly(handleableAnomalyDTO HandleableAnomalyDTO) sensoranomalyhandler.HandleableAnomaly {
	return sensoranomalyhandler.HandleableAnomaly{
		Anomaly: sensoranomalyhandler.Anomaly{
			AnomalyType: handleableAnomalyDTO.AnomalyType,
		},
	}
}

func mapControllCommandsDtoToControllCommandsModel(controllCommandsDTO []ControllCommandDTO) []models.ControllCommandModel {
	controllCommandsModel := []models.ControllCommandModel{}
	for _, controllCommandDTO := range controllCommandsDTO {
		controllCommandsModel = append(controllCommandsModel, mapControllCommandDtoToControllCommandModel(controllCommandDTO))
	}
	return controllCommandsModel
}

func mapControllCommandDtoToControllCommandModel(sensedPropertyDTO ControllCommandDTO) models.ControllCommandModel {
	return models.ControllCommandModel{
		ServiceDefinition: serviceModels.ServiceDefinition{
			ServiceDefinition: sensedPropertyDTO.ServiceDefinition,
			ServiceUri:        sensedPropertyDTO.ServiceUri,
		},
	}
}

func mapSensedPropertiesDtoToSensedPropertiesModel(sensedPropertiesDTO []SensedPropertyDTO) []models.SensedPropertyModel {
	sensedPropertyModel := []models.SensedPropertyModel{}
	for _, sensesensedPropertyDTO := range sensedPropertiesDTO {
		sensedPropertyModel = append(sensedPropertyModel, mapSensedPropertyDtoToSensedPropertyModel(sensesensedPropertyDTO))
	}
	return sensedPropertyModel
}

func mapSensedPropertyDtoToSensedPropertyModel(sensedPropertyDTO SensedPropertyDTO) models.SensedPropertyModel {
	return models.SensedPropertyModel{
		ServiceDefinition: serviceModels.ServiceDefinition{
			ServiceDefinition: sensedPropertyDTO.ServiceDefinition,
			ServiceUri:        sensedPropertyDTO.ServiceUri,
		},
		SensorEndpointMode: models.SensorEndpointMode(sensedPropertyDTO.SensorEndpointMode),
		IntervalTime:       sensedPropertyDTO.IntervalTime,
	}
}
