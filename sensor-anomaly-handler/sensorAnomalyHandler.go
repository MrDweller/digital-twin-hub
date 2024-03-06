package sensoranomalyhandler

import (
	"net/http"

	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	"github.com/MrDweller/service-registry-connection/models"
	"github.com/gin-gonic/gin"
)

type Service interface {
	HandleAnomaly(anomaly Anomaly) error
}

type SensorAnomalyHandler struct {
	router              *gin.Engine
	handleableAnomalies []HandleableAnomaly
	service             Service
	handlingSystem      models.SystemDefinition
}

func InitAnomalyHandler(handleableAnomalies []HandleableAnomaly, router *gin.Engine, SystemName string) []additionalservice.AdditionalService {
	additionalServices := []additionalservice.AdditionalService{}
	handlingSystem := models.SystemDefinition{
		Address:    "localhost",
		Port:       5672,
		SystemName: SystemName,
	}
	for _, handleableAnomaly := range handleableAnomalies {
		handleableAnomaly.setAnomalyHandlingSystem(handlingSystem)
		additionalServices = append(additionalServices, &handleableAnomaly)
		additionalServices = append(additionalServices, &NotifyAnomaly{
			Anomaly: handleableAnomaly.Anomaly,
		})
	}
	sensorAnomalyHandler := newSensorAnomalyHandler(router, handleableAnomalies, handlingSystem)
	sensorAnomalyHandler.SetupEndpoints()

	return additionalServices
}

func newSensorAnomalyHandler(router *gin.Engine, handleableAnomalies []HandleableAnomaly, handlingSystem models.SystemDefinition) SensorAnomalyHandler {
	return SensorAnomalyHandler{
		router:              router,
		handleableAnomalies: handleableAnomalies,
		service: RabbitmqAnomalyHandlerService{
			rabbitmqAddress: "localhost",
			rabbitmqPort:    5672,
		},
		handlingSystem: handlingSystem,
	}

}

func (sensorAnomalyHandler SensorAnomalyHandler) SetupEndpoints() {
	for _, handleableAnomaly := range sensorAnomalyHandler.handleableAnomalies {

		sensorAnomalyHandler.router.POST(
			handleableAnomaly.GetService().ServiceDefinition.ServiceUri,
			func(c *gin.Context) {

				err := sensorAnomalyHandler.service.HandleAnomaly(handleableAnomaly.Anomaly)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.Status(http.StatusOK)
			},
		)

	}
}
