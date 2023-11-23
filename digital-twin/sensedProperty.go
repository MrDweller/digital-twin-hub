package digitaltwin

import (
	"net/http"

	"github.com/MrDweller/digital-twin-hub/models"
	physicaltwinconnection "github.com/MrDweller/digital-twin-hub/physical-twin-connection"
	"github.com/gin-gonic/gin"
)

func AddSensorEnpoint(router *gin.Engine, sensedPropertyModel models.SensedPropertyModel, connection physicaltwinconnection.Connection) {
	router.GET(sensedPropertyModel.ServiceUri, func(c *gin.Context) {
		response, err := connection.HandleSensorRequest(sensedPropertyModel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, response)
	})
}
