package manufacturer

import (
	digitaltwin "github.com/MrDweller/digital-twin-hub/digital-twin"
	digitaltwinmodels "github.com/MrDweller/digital-twin-hub/digital-twin-models"
)

type Service struct{}

func NewService() Service {
	service := Service{}
	return service
}

func (service Service) CreateDigitalTwin(digitalTwinModel digitaltwin.DigitalTwinModel) (*digitaltwinmodels.SystemDefinition, error) {
	systemDefinition, err := digitalTwinModel.StartDigitalTwin()
	if err != nil {
		return nil, err
	}

	return systemDefinition, nil
}
