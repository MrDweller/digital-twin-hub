package manufacturer

import (
	"os"
	"strconv"

	digitaltwin "github.com/MrDweller/digital-twin-hub/digital-twin"
	digitaltwinmodels "github.com/MrDweller/digital-twin-hub/digital-twin-models"
	digitaltwinregistry "github.com/MrDweller/digital-twin-hub/digital-twin-registry"
)

type Service struct {
	digitalTwinRegistryConnection digitaltwinregistry.DigitalTwinRegistryConnection
}

func NewService() (*Service, error) {
	srPort, err := strconv.Atoi(os.Getenv("SR_PORT"))
	if err != nil {
		return nil, err
	}
	digitalTwinRegistryConnection, err := digitaltwinregistry.NewConnection(digitaltwinregistry.DigitalTwinRegistry{
		Address: os.Getenv("SR_ADDRESS"),
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

func (service Service) CreateDigitalTwin(digitalTwinModel digitaltwinmodels.DigitalTwinModel) (*digitaltwinmodels.SystemDefinition, error) {
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
