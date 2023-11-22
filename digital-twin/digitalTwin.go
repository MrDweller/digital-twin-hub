package digitaltwin

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	digitaltwinmodels "github.com/MrDweller/digital-twin-hub/digital-twin-models"
	digitaltwinregistry "github.com/MrDweller/digital-twin-hub/digital-twin-registry"
	physicaltwinconnection "github.com/MrDweller/digital-twin-hub/physical-twin-connection"

	"github.com/gin-gonic/gin"
)

type DigitalTwin struct {
	DigitalTwinModel    digitaltwinmodels.DigitalTwinModel
	digitalTwinRegistry digitaltwinregistry.DigitalTwinRegistryConnection
	systemDefinition    digitaltwinmodels.SystemDefinition
	listener            net.Listener
}

func NewDigitalTwin(digitalTwinModel digitaltwinmodels.DigitalTwinModel, digitalTwinRegistryConnection digitaltwinregistry.DigitalTwinRegistryConnection) (*DigitalTwin, error) {
	url := fmt.Sprintf("%s:0", os.Getenv("ADDRESS"))
	listener, err := net.Listen("tcp", url)
	if err != nil {
		return nil, err
	}

	address, stringPort, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(stringPort)
	if err != nil {
		return nil, err
	}
	systemDefinition := digitaltwinmodels.SystemDefinition{
		Address:    address,
		Port:       port,
		SystemName: os.Getenv("SYSTEM_NAME"),
	}

	return &DigitalTwin{
		DigitalTwinModel:    digitalTwinModel,
		digitalTwinRegistry: digitalTwinRegistryConnection,
		systemDefinition:    systemDefinition,
		listener:            listener,
	}, nil
}

func (digitalTwin *DigitalTwin) StartDigitalTwin() (*digitaltwinmodels.SystemDefinition, error) {
	router := gin.New()

	connection, err := physicaltwinconnection.NewConnection(digitalTwin.DigitalTwinModel.PhysicalTwinConnectionModel)
	if err != nil {
		return nil, err
	}

	for _, sensedPropertyModel := range digitalTwin.DigitalTwinModel.SensedProperties {
		AddSensorEnpoint(router, sensedPropertyModel, connection)
	}
	for _, controlCommandModel := range digitalTwin.DigitalTwinModel.ControlCommands {
		AddCommandEnpoint(router, controlCommandModel, connection)
	}

	err = digitalTwin.digitalTwinRegistry.RegisterDigitalTwin(digitalTwin.DigitalTwinModel, digitalTwin.systemDefinition)
	if err != nil {
		return nil, err
	}

	// Start the digital twin's rest api
	go http.Serve(digitalTwin.listener, router)

	return &digitalTwin.systemDefinition, nil
}
