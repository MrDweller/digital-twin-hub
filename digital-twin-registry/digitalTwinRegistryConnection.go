package digitaltwinregistry

import (
	"errors"
	"fmt"
	"os"

	"github.com/MrDweller/digital-twin-hub/models"
	serviceModels "github.com/MrDweller/service-registry-connection/models"
)

type DigitalTwinRegistryConnection interface {
	connect() error
	RegisterDigitalTwin(digitalTwinModel models.DigitalTwinModel, systemDefinition serviceModels.SystemDefinition) error
	UnRegisterDigitalTwin(digitalTwinModel models.DigitalTwinModel, systemDefinition serviceModels.SystemDefinition) error
}

type DigitalTwinRegistryImplementationType string

func NewConnection(digitalTwinRegistry DigitalTwinRegistry) (DigitalTwinRegistryConnection, error) {
	var digitalTwinRegistryImplementationType DigitalTwinRegistryImplementationType
	digitalTwinRegistryImplementationType = DigitalTwinRegistryImplementationType(os.Getenv("DIGITAL_TWIN_REGISTRY_IMPLEMENTATION"))

	var digitalTwinRegistryConnection DigitalTwinRegistryConnection

	switch digitalTwinRegistryImplementationType {
	case DIGITAL_TWIN_REGISTRY_ARROWHEAD_4_6_1:
		digitalTwinRegistryConnection = NewDigitalTwinRegistryArrowhead_4_6_1(digitalTwinRegistry)
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
