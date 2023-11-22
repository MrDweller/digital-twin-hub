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
}

func NewDigitalTwin(digitalTwinModel digitaltwinmodels.DigitalTwinModel, digitalTwinRegistryConnection digitaltwinregistry.DigitalTwinRegistryConnection) *DigitalTwin {
	return &DigitalTwin{
		DigitalTwinModel:    digitalTwinModel,
		digitalTwinRegistry: digitalTwinRegistryConnection,
	}
}

func (digitalTwin *DigitalTwin) StartDigitalTwin() (*digitaltwinmodels.SystemDefinition, error) {
	router := gin.New()

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

	systemDefinition := digitaltwinmodels.SystemDefinition{
		Address:    address,
		Port:       port,
		SystemName: os.Getenv("SYSTEM_NAME"),
	}

	err = digitalTwin.digitalTwinRegistry.RegisterDigitalTwin(digitalTwin.DigitalTwinModel, systemDefinition)
	if err != nil {
		return nil, err
	}

	// Start the digital twin's rest api
	go http.Serve(listener, router)

	return &systemDefinition, nil
}
