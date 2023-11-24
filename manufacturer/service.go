package manufacturer

import (
	"os"
	"strconv"

	digitaltwin "github.com/MrDweller/digital-twin-hub/digital-twin"
	digitaltwinregistry "github.com/MrDweller/digital-twin-hub/digital-twin-registry"
	"github.com/MrDweller/digital-twin-hub/models"
	serviceModels "github.com/MrDweller/service-registry-connection/models"
)

type Service struct {
	digitalTwinRegistryConnection digitaltwinregistry.DigitalTwinRegistryConnection
}

func NewService() (*Service, error) {
	srPort, err := strconv.Atoi(os.Getenv("DIGITAL_TWIN_REGISTRY_PORT"))
	if err != nil {
		return nil, err
	}
	digitalTwinRegistryConnection, err := digitaltwinregistry.NewConnection(digitaltwinregistry.DigitalTwinRegistry{
		Address: os.Getenv("DIGITAL_TWIN_REGISTRY_ADDRESS"),
		Port:    srPort,
	})
	if err != nil {
		return nil, err
	}
	service := &Service{
		digitalTwinRegistryConnection: digitalTwinRegistryConnection,
	}
	return service, nil
}

func (service Service) CreateDigitalTwin(digitalTwinModel models.DigitalTwinModel) (*serviceModels.SystemDefinition, error) {
	digitalTwin, err := digitaltwin.NewDigitalTwin(digitalTwinModel, service.digitalTwinRegistryConnection)
	if err != nil {
		return nil, err
	}

	systemDefinition, err := digitalTwin.StartDigitalTwin()
	if err != nil {
		return nil, err
	}

	return systemDefinition, nil
}
