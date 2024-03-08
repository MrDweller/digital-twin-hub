package sensoranomalyhandler

import (
	"fmt"
	"net/http"

	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	"github.com/MrDweller/service-registry-connection/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SensorAnomalyHandler struct {
	router                         *gin.Engine
	notifyToHandleableAnomaliesMap map[*NotifyAnomaly]HandleableAnomaly
	anomalyHandlerService          AnomalyHandlerService
}

func InitAnomalyHandler(anomalies []Anomaly, router *gin.Engine, SystemName string) []additionalservice.AdditionalService {
	additionalServices := []additionalservice.AdditionalService{}
	handlingSystem := models.SystemDefinition{
		Address:    "localhost",
		Port:       5672,
		SystemName: SystemName,
	}

	notifyToHandleableAnomaliesMap := map[*NotifyAnomaly]HandleableAnomaly{}
	for _, anomaly := range anomalies {

		handleableAnomaly := &RabbitMQHandleableAnomaly{
			HandleableAnomalyBase: HandleableAnomalyBase{
				anomalyHandlingSystem: handlingSystem,
				Anomaly:               anomaly,
			},
			id: uuid.New(),
		}

		additionalServices = append(additionalServices, handleableAnomaly)

		notifyAnomaly := &NotifyAnomaly{
			Anomaly: anomaly,
		}
		additionalServices = append(additionalServices, notifyAnomaly)

		notifyToHandleableAnomaliesMap[notifyAnomaly] = handleableAnomaly
		fmt.Println(handleableAnomaly.GetService())
	}
	sensorAnomalyHandler := newSensorAnomalyHandler(router, notifyToHandleableAnomaliesMap, RabbitmqAnomalyHandlerService{
		rabbitmqAddress: handlingSystem.Address,
		rabbitmqPort:    handlingSystem.Port,
	})
	sensorAnomalyHandler.SetupEndpoints()

	return additionalServices
}

func newSensorAnomalyHandler(router *gin.Engine, notifyToHandleableAnomaliesMap map[*NotifyAnomaly]HandleableAnomaly, anomalyHandlerService AnomalyHandlerService) SensorAnomalyHandler {
	return SensorAnomalyHandler{
		router:                         router,
		notifyToHandleableAnomaliesMap: notifyToHandleableAnomaliesMap,
		anomalyHandlerService:          anomalyHandlerService,
	}

}

func (sensorAnomalyHandler *SensorAnomalyHandler) SetupEndpoints() {
	for notifyAnomaly, handleableAnomaly := range sensorAnomalyHandler.notifyToHandleableAnomaliesMap {
		sensorAnomalyHandler.addNotifyEndpoint(notifyAnomaly, handleableAnomaly)

	}
}

func (sensorAnomalyHandler *SensorAnomalyHandler) addNotifyEndpoint(notifyAnomaly *NotifyAnomaly, handleableAnomaly HandleableAnomaly) {
	sensorAnomalyHandler.router.POST(
		notifyAnomaly.GetService().ServiceDefinition.ServiceUri,
		func(c *gin.Context) {
			err := sensorAnomalyHandler.anomalyHandlerService.HandleAnomaly(handleableAnomaly)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.Status(http.StatusOK)
		},
	)
}
