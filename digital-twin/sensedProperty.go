package digitaltwin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MrDweller/digital-twin-hub/database"
	"github.com/MrDweller/digital-twin-hub/models"
	physicaltwinconnection "github.com/MrDweller/digital-twin-hub/physical-twin-connection"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddSensorEnpoint(router *gin.Engine, digitalTwinId uuid.UUID, sensedPropertyModel models.SensedPropertyModel, connection physicaltwinconnection.Connection) error {
	var sensorEndpoint SensorEndpoint
	switch sensedPropertyModel.SensorEndpointMode {
	case models.IMMEDIATE_RETRIEVAL:
		sensorEndpoint = newImmediateSensorEndpoint(digitalTwinId, sensedPropertyModel, connection)
	case models.INTERVAL_RETRIEVAL:
		sensorEndpoint = newIntervalSensorEndpoint(digitalTwinId, sensedPropertyModel, connection)
	default:
		errorString := fmt.Sprintf("the sensor endpoint mode %s has no implementation", sensedPropertyModel.SensorEndpointMode)
		return errors.New(errorString)
	}

	router.GET(sensedPropertyModel.ServiceUri, func(c *gin.Context) {
		sensedData, err := sensorEndpoint.HandleRequest()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, sensedData)
	})

	return nil
}

type SensorEndpoint interface {
	HandleRequest() (map[string]any, error)
	retriveSensorData() (map[string]any, error)
}

type ImmediateSensorEndpoint struct {
	SensedProperty
}

func newImmediateSensorEndpoint(digitalTwinId uuid.UUID, sensedPropertyModel models.SensedPropertyModel, connection physicaltwinconnection.Connection) ImmediateSensorEndpoint {
	sensorEndpoint := ImmediateSensorEndpoint{
		SensedProperty: SensedProperty{
			DigitalTwinId:       digitalTwinId,
			SensedPropertyModel: sensedPropertyModel,
			connection:          connection,
		},
	}
	return sensorEndpoint
}

func (immediateSensorEndpoint ImmediateSensorEndpoint) HandleRequest() (map[string]any, error) {
	sensedData, err := immediateSensorEndpoint.retriveSensorData()
	if err != nil {
		log.Println(err)
	}
	err = immediateSensorEndpoint.saveSensorData(sensedData)
	if err != nil {
		log.Println(err)
	}
	return immediateSensorEndpoint.getStoredSensorData()
}

type IntervalSensorEndpoint struct {
	SensedProperty
}

func newIntervalSensorEndpoint(digitalTwinId uuid.UUID, sensedPropertyModel models.SensedPropertyModel, connection physicaltwinconnection.Connection) IntervalSensorEndpoint {
	sensorEndpoint := IntervalSensorEndpoint{
		SensedProperty: SensedProperty{
			DigitalTwinId:       digitalTwinId,
			SensedPropertyModel: sensedPropertyModel,
			connection:          connection,
		},
	}

	ticker := time.NewTicker(time.Duration(sensedPropertyModel.IntervalTime) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				sensedData, err := sensorEndpoint.retriveSensorData()
				if err != nil {
					log.Println(err)
					continue
				}
				err = sensorEndpoint.saveSensorData(sensedData)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()
	return sensorEndpoint
}

func (intervalSensorEndpoint IntervalSensorEndpoint) HandleRequest() (map[string]any, error) {
	return intervalSensorEndpoint.getStoredSensorData()
}

type SensedProperty struct {
	DigitalTwinId uuid.UUID
	models.SensedPropertyModel
	connection physicaltwinconnection.Connection
}

func (sensedProperty SensedProperty) retriveSensorData() (map[string]any, error) {
	response, err := sensedProperty.connection.HandleSensorRequest(sensedProperty.SensedPropertyModel)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (sensedProperty SensedProperty) getStoredSensorData() (map[string]any, error) {
	filter := bson.M{
		"digitaltwinid":     sensedProperty.DigitalTwinId,
		"servicedefinition": sensedProperty.ServiceDefinition.ServiceDefinition,
	}
	result := database.SensorData.FindOne(context.Background(), filter)

	var sensorData SensorData
	if err := result.Decode(&sensorData); err != nil {
		return nil, err
	}
	return sensorData.SensedData, nil
}

func (sensedProperty SensedProperty) saveSensorData(sensedData map[string]any) error {
	filter := bson.M{
		"digitaltwinid":     sensedProperty.DigitalTwinId,
		"servicedefinition": sensedProperty.ServiceDefinition.ServiceDefinition,
	}
	// Try to find the object by ID
	result := database.SensorData.FindOne(context.Background(), filter)

	if result.Err() == nil {
		// Object exists, update the data property
		update := bson.M{"$set": bson.M{"senseddata": sensedData}}
		_, err := database.SensorData.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return err
		}
		log.Println("Object updated successfully")
	} else if result.Err() == mongo.ErrNoDocuments {
		// Object does not exist, insert a new one
		_, err := database.SensorData.InsertOne(context.Background(), SensorData{
			DigitalTwinId:     sensedProperty.DigitalTwinId,
			ServiceDefinition: sensedProperty.ServiceDefinition.ServiceDefinition,
			SensedData:        sensedData,
		})
		if err != nil {
			return err
		}
		log.Println("Object inserted successfully")
	} else {
		// Handle other errors
		return errors.New("could not save sensor data to the database")
	}
	return nil

}

type SensorData struct {
	DigitalTwinId     uuid.UUID
	ServiceDefinition string
	SensedData        map[string]any
}
