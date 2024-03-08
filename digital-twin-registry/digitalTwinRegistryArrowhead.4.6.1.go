package digitaltwinregistry

import (
	"log"

	"github.com/MrDweller/digital-twin-hub/models"
	serviceModels "github.com/MrDweller/service-registry-connection/models"
	serviceregistry "github.com/MrDweller/service-registry-connection/service-registry"
)

const DIGITAL_TWIN_REGISTRY_ARROWHEAD_4_6_1 DigitalTwinRegistryImplementationType = "digital-twin-registry-arrowhead-4.6.1"

type DigitalTwinRegistryArrowhead_4_6_1 struct {
	DigitalTwinRegistry
	serviceregistry.ServiceRegistryConnection
}

func NewDigitalTwinRegistryArrowhead_4_6_1(digitalTwinRegistry DigitalTwinRegistry) DigitalTwinRegistryArrowhead_4_6_1 {
	return DigitalTwinRegistryArrowhead_4_6_1{
		DigitalTwinRegistry: digitalTwinRegistry,
		ServiceRegistryConnection: serviceregistry.ServiceRegistryArrowhead_4_6_1{
			ServiceRegistry: serviceregistry.ServiceRegistry{
				Address: digitalTwinRegistry.Address,
				Port:    digitalTwinRegistry.Port,
			},
			CertificateInfo: digitalTwinRegistry.CertificateInfo,
		},
	}
}

func (digitalTwinRegistry DigitalTwinRegistryArrowhead_4_6_1) connect() error {
	err := digitalTwinRegistry.ServiceRegistryConnection.Connect()
	if err != nil {
		return err
	}

	return nil

}

func (digitalTwinRegistry DigitalTwinRegistryArrowhead_4_6_1) RegisterDigitalTwin(digitalTwinModel models.DigitalTwinModel, systemDefinition serviceModels.SystemDefinition) error {
	_, err := digitalTwinRegistry.RegisterSystem(systemDefinition)
	if err != nil {
		return err
	}

	for _, sensedPropertyModel := range digitalTwinModel.SensedProperties {
		_, err := digitalTwinRegistry.RegisterService(
			sensedPropertyModel.ServiceDefinition,
			[]string{
				"HTTP-SECURE-JSON",
			},
			map[string]string{},
			systemDefinition,
		)
		if err != nil {
			return err
		}

	}

	for _, controlCommandModel := range digitalTwinModel.ControlCommands {
		_, err := digitalTwinRegistry.RegisterService(
			controlCommandModel.ServiceDefinition,
			[]string{
				"HTTP-SECURE-JSON",
			},
			map[string]string{},
			systemDefinition,
		)
		if err != nil {
			return err
		}

	}

	for _, additionalServiceModel := range digitalTwinModel.AdditionalServiceModels {
		log.Println(additionalServiceModel)
		var localSystemDefinition serviceModels.SystemDefinition
		if additionalServiceModel.HasExternalHost {
			localSystemDefinition = additionalServiceModel.SystemDefinition
		} else {
			localSystemDefinition = systemDefinition
		}

		_, err := digitalTwinRegistry.RegisterService(
			additionalServiceModel.ServiceDefinition,
			additionalServiceModel.Interfaces,
			additionalServiceModel.Metadata,
			localSystemDefinition,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (digitalTwinRegistry DigitalTwinRegistryArrowhead_4_6_1) UnRegisterDigitalTwin(digitalTwinModel models.DigitalTwinModel, systemDefinition serviceModels.SystemDefinition) error {
	for _, sensedPropertyModel := range digitalTwinModel.SensedProperties {
		err := digitalTwinRegistry.UnRegisterService(sensedPropertyModel.ServiceDefinition, systemDefinition)
		if err != nil {
			return err
		}

	}

	for _, controlCommandModel := range digitalTwinModel.ControlCommands {
		err := digitalTwinRegistry.UnRegisterService(controlCommandModel.ServiceDefinition, systemDefinition)
		if err != nil {
			return err
		}

	}

	for _, additionalServiceModel := range digitalTwinModel.AdditionalServiceModels {
		if additionalServiceModel.HasExternalHost {
			digitalTwinRegistry.UnRegisterSystem(additionalServiceModel.SystemDefinition)
		} else {
			err := digitalTwinRegistry.UnRegisterService(
				additionalServiceModel.ServiceDefinition,
				systemDefinition,
			)
			if err != nil {
				return err
			}
		}
	}

	return digitalTwinRegistry.UnRegisterSystem(systemDefinition)
}
