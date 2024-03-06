package manufacturer

import (
	"context"
	"os"
	"strconv"

	"github.com/MrDweller/digital-twin-hub/database"
	digitaltwin "github.com/MrDweller/digital-twin-hub/digital-twin"
	digitaltwinregistry "github.com/MrDweller/digital-twin-hub/digital-twin-registry"
	"github.com/MrDweller/digital-twin-hub/models"
	serviceModels "github.com/MrDweller/service-registry-connection/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (service Service) CreateDigitalTwin(digitalTwinModel models.DigitalTwinModel, digitalTwinId uuid.UUID, router *gin.Engine) (*serviceModels.SystemDefinition, error) {

	digitalTwin, err := digitaltwin.NewDigitalTwin(digitalTwinModel, service.digitalTwinRegistryConnection, digitalTwinId, router)
	if err != nil {
		return nil, err
	}

	database.DigitalTwin.InsertOne(context.Background(), digitalTwin)

	systemDefinition, err := digitalTwin.StartDigitalTwin()
	if err != nil {
		return nil, err
	}

	return systemDefinition, nil
}

func (service Service) DeleteDigitalTwin(address string, port int) error {
	filter := bson.M{
		"systemdefinition.address": address,
		"systemdefinition.port":    port,
	}
	digitalTwins, err := service.getSavedDigitalTwins(filter)
	database.DigitalTwin.DeleteMany(context.Background(), filter)

	if err != nil {
		return err
	}

	for _, didigitalTwin := range digitalTwins {
		database.SensorData.DeleteMany(context.Background(), bson.M{
			"digitaltwinid": didigitalTwin.DigitalTwinId,
		})
		err := service.digitalTwinRegistryConnection.UnRegisterDigitalTwin(didigitalTwin.DigitalTwinModel, didigitalTwin.SystemDefinition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service Service) startAllSavedDigitalTwins() error {
	digitalTwins, err := service.getSavedDigitalTwins(bson.M{})
	if err != nil {
		return err
	}

	for _, digitalTwin := range digitalTwins {
		filter := bson.M{
			"systemdefinition.address": digitalTwin.SystemDefinition.Address,
			"systemdefinition.port":    digitalTwin.SystemDefinition.Port,
		}
		_, err := database.DigitalTwin.DeleteMany(context.Background(), filter)
		if err != nil {
			return err
		}
	}

	for _, didigitalTwin := range digitalTwins {
		router := gin.New()
		_, err := service.CreateDigitalTwin(didigitalTwin.DigitalTwinModel, didigitalTwin.DigitalTwinId, router)
		if err != nil {
			return err
		}
	}
	return nil

}

func (service Service) stopAllSavedDigitalTwins() error {
	digitalTwins, err := service.getSavedDigitalTwins(bson.M{})
	if err != nil {
		return err
	}

	for _, digitalTwin := range digitalTwins {
		err := service.digitalTwinRegistryConnection.UnRegisterDigitalTwin(digitalTwin.DigitalTwinModel, digitalTwin.SystemDefinition)
		if err != nil {
			return err
		}
	}
	return nil

}

func (service Service) getSavedDigitalTwins(filter interface{}, opts ...*options.FindOptions) ([]digitaltwin.DigitalTwin, error) {
	result, err := database.DigitalTwin.Find(context.Background(), filter, opts...)
	if err != nil {
		return nil, err
	}

	var digitalTwins []digitaltwin.DigitalTwin
	if err := result.All(context.Background(), &digitalTwins); err != nil {
		return nil, err
	}

	return digitalTwins, nil
}
