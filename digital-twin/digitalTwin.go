package digitaltwin

import (
	"fmt"
	"net"
	"os"
	"strconv"

	digitaltwinregistry "github.com/MrDweller/digital-twin-hub/digital-twin-registry"
	httpserver "github.com/MrDweller/digital-twin-hub/http-server"
	"github.com/MrDweller/digital-twin-hub/models"
	physicaltwinconnection "github.com/MrDweller/digital-twin-hub/physical-twin-connection"
	serviceModels "github.com/MrDweller/service-registry-connection/models"

	"github.com/gin-gonic/gin"
)

type DigitalTwin struct {
	DigitalTwinModel    models.DigitalTwinModel
	digitalTwinRegistry digitaltwinregistry.DigitalTwinRegistryConnection
	systemDefinition    serviceModels.SystemDefinition
	server              httpserver.Server
}

func NewDigitalTwin(digitalTwinModel models.DigitalTwinModel, digitalTwinRegistryConnection digitaltwinregistry.DigitalTwinRegistryConnection) (*DigitalTwin, error) {
	url := fmt.Sprintf("%s:0", os.Getenv("ADDRESS"))

	connection, err := physicaltwinconnection.NewConnection(digitalTwinModel.PhysicalTwinConnectionModel)
	if err != nil {
		return nil, err
	}
	router := gin.New()
	for _, sensedPropertyModel := range digitalTwinModel.SensedProperties {
		AddSensorEnpoint(router, sensedPropertyModel, connection)
	}
	for _, controlCommandModel := range digitalTwinModel.ControlCommands {
		AddCommandEnpoint(router, controlCommandModel, connection)
	}
	server, err := httpserver.NewServer(url, router)
	if err != nil {
		return nil, err
	}

	address, stringPort, err := net.SplitHostPort(server.Addr())
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(stringPort)
	if err != nil {
		return nil, err
	}
	systemDefinition := serviceModels.SystemDefinition{
		Address:            address,
		Port:               port,
		SystemName:         os.Getenv("SYSTEM_NAME"),
		AuthenticationInfo: os.Getenv("AUTHENTICATION_INFO"),
	}

	return &DigitalTwin{
		DigitalTwinModel:    digitalTwinModel,
		digitalTwinRegistry: digitalTwinRegistryConnection,
		systemDefinition:    systemDefinition,
		server:              server,
	}, nil
}

func (digitalTwin *DigitalTwin) StartDigitalTwin() (*serviceModels.SystemDefinition, error) {

	err := digitalTwin.digitalTwinRegistry.RegisterDigitalTwin(digitalTwin.DigitalTwinModel, digitalTwin.systemDefinition)
	if err != nil {
		return nil, err
	}

	go digitalTwin.server.StartServer()

	return &digitalTwin.systemDefinition, nil
}
