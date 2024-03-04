package sensoranomalyhandler

import "github.com/gin-gonic/gin"

type Anomaly struct {
	AnomalyType string `json:"anomalyType"`
}

type Service interface {
	HandleAnomaly(anomaly Anomaly) error
}

type SensorAnomalyHandler struct {
	router  *gin.Engine
	service Service
}

func NewSensorAnomalyHandler(router *gin.Engine) SensorAnomalyHandler {
	return SensorAnomalyHandler{
		router:  router,
		service: RabbitmqAnomalyHandlerService{},
	}

}

func (sensorAnomalyHandler SensorAnomalyHandler) SetupEndpoints() {
	controller := NewController(sensorAnomalyHandler.service)
	sensorAnomalyHandler.router.POST("/sensor-anomaly", controller.SensorAnomaly)
}
