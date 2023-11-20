package physicaltwinconnection

import (
	"errors"
	"fmt"

	digitaltwinmodels "github.com/MrDweller/digital-twin-hub/digital-twin-models"
	"github.com/mitchellh/mapstructure"
)

type Connection interface {
	connect() error
	HandleControllCommand(serviceDefinition digitaltwinmodels.ControllCommandModel, commands any) (string, error)
	HandleSensorRequest(serviceDefinition digitaltwinmodels.SensedPropertyModel) (string, error)
}

type PhysicalTwinConnectionModel struct {
	ConnectionType  PhysicalTwinConnectionType `json:"connectionType"`
	ConnectionModel map[string]any             `json:"connectionModel"`
}

func NewConnection(physicalTwinConnection PhysicalTwinConnectionModel) (Connection, error) {
	var connection Connection

	switch physicalTwinConnection.ConnectionType {
	case SIMPLE_COAP:
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
