package digitaltwin

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	digitaltwinmodels "github.com/MrDweller/digital-twin-hub/digital-twin-models"
	physicaltwinconnection "github.com/MrDweller/digital-twin-hub/physical-twin-connection"

	"github.com/gin-gonic/gin"
)

type DigitalTwinModel struct {
	SensedProperties []digitaltwinmodels.SensedPropertyModel  `json:"sensedProperties"`
	ControlCommands  []digitaltwinmodels.ControllCommandModel `json:"controlCommands"`

	PhysicalTwinConnectionModel physicaltwinconnection.PhysicalTwinConnectionModel `json:"physicalTwinConnection"`
}

func (digitalTwinModel *DigitalTwinModel) StartDigitalTwin() (*digitaltwinmodels.SystemDefinition, error) {
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

	connection, err := physicaltwinconnection.NewConnection(digitalTwinModel.PhysicalTwinConnectionModel)
	if err != nil {
		return nil, err
	}

	for _, sensedPropertyModel := range digitalTwinModel.SensedProperties {
		AddSensorEnpoint(router, sensedPropertyModel, connection)
	}
	for _, controlCommandModel := range digitalTwinModel.ControlCommands {
		AddCommandEnpoint(router, controlCommandModel, connection)
	}

	// Start the digital twin's rest api
	go http.Serve(listener, router)

	return &digitaltwinmodels.SystemDefinition{
		Address: address,
		Port:    port,
	}, nil
}
