package digitaltwinregistry

import (
	"errors"
	"fmt"
	"os"

	digitaltwinmodels "github.com/MrDweller/digital-twin-hub/digital-twin-models"
)

type DigitalTwinRegistryConnection interface {
	connect() error
	RegisterDigitalTwin(digitalTwinModel digitaltwinmodels.DigitalTwinModel, systemDefinition digitaltwinmodels.SystemDefinition) error
	UnRegisterDigitalTwin(digitalTwinModel digitaltwinmodels.DigitalTwinModel, systemDefinition digitaltwinmodels.SystemDefinition) error
}

type DigitalTwinRegistryImplementationType string

func NewConnection(digitalTwinRegistry DigitalTwinRegistry) (DigitalTwinRegistryConnection, error) {
	var digitalTwinRegistryImplementationType DigitalTwinRegistryImplementationType
	digitalTwinRegistryImplementationType = DigitalTwinRegistryImplementationType(os.Getenv("DIGITAL_TWIN_REGISTRY_IMPLEMENTATION"))

	var digitalTwinRegistryConnection DigitalTwinRegistryConnection

	switch digitalTwinRegistryImplementationType {
	case SERVICE_REGISTRY_ARROWHEAD_4_6_1:
		digitalTwinRegistryConnection = ServiceRegistryArrowhead_4_6_1{
			DigitalTwinRegistry: digitalTwinRegistry,
		}
		break
	default:
		errorString := fmt.Sprintf("the digital twin registry %s has no implementation", digitalTwinRegistryImplementationType)
		return nil, errors.New(errorString)
	}

	err := digitalTwinRegistryConnection.connect()
	if err != nil {
		return nil, err
	}

	return digitalTwinRegistryConnection, nil
}
