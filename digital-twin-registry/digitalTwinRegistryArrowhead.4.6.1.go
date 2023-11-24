package digitaltwinregistry

import (
	"os"

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
			CertificateInfo: serviceModels.CertificateInfo{
				CertFilePath: os.Getenv("CERT_FILE_PATH"),
				KeyFilePath:  os.Getenv("KEY_FILE_PATH"),
				Truststore:   os.Getenv("TRUSTSTORE_FILE_PATH"),
			},
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
	for _, sensedPropertyModel := range digitalTwinModel.SensedProperties {
		_, err := digitalTwinRegistry.RegisterService(sensedPropertyModel.ServiceDefinition, systemDefinition)
		if err != nil {
			return err
		}

	}

	for _, controlCommandModel := range digitalTwinModel.ControlCommands {
		_, err := digitalTwinRegistry.RegisterService(controlCommandModel.ServiceDefinition, systemDefinition)
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
	return digitalTwinRegistry.UnRegisterSystem(systemDefinition)
}
