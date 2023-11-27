package physicaltwinconnection

import (
	"errors"
	"fmt"

	"github.com/MrDweller/digital-twin-hub/models"
	physicaltwinmodels "github.com/MrDweller/digital-twin-hub/models"
	"github.com/mitchellh/mapstructure"
)

type Connection interface {
	connect() error
	HandleControllCommand(serviceDefinition models.ControllCommandModel, commands any) (map[string]any, error)
	HandleSensorRequest(serviceDefinition models.SensedPropertyModel) (map[string]any, error)
}

func NewConnection(physicalTwinConnection physicaltwinmodels.PhysicalTwinConnectionModel) (Connection, error) {
	var connection Connection

	switch physicalTwinConnection.ConnectionType {
	case physicaltwinmodels.SIMPLE_COAP:
		var simpleCoapConnectionModel SimpleCoapConnectionModel
		mapstructure.Decode(physicalTwinConnection.ConnectionModel, &simpleCoapConnectionModel)
		connection = simpleCoapConnectionModel
		break

	default:
		errorString := fmt.Sprintf("the physical twin connection type %s has no implementation", physicalTwinConnection.ConnectionType)
		return nil, errors.New(errorString)
	}

	err := connection.connect()
	if err != nil {
		return nil, err
	}

	return connection, nil
}
