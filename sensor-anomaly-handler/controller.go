package sensoranomalyhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

// Sensor anomaly
// @Summary      Notyfies about anomalies.
// @Description  Notyfies the digital twin hub of a detected anomaly in the sensor data.
// @Tags         Management
// @Param        anomaly 				body       	Anomaly  	true  "{anomalyType: "STUCK"}"
// @Success      200
// @Router       /sensor-anomaly [post]
func (controller *Controller) SensorAnomaly(c *gin.Context) {
	var anomaly Anomaly
	err := c.BindJSON(&anomaly)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = controller.service.HandleAnomaly(anomaly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
