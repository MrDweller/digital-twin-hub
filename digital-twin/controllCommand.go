package digitaltwin

import (
	"net/http"

	"github.com/MrDweller/digital-twin-hub/models"
	physicaltwinconnection "github.com/MrDweller/digital-twin-hub/physical-twin-connection"
	"github.com/gin-gonic/gin"
)

func AddCommandEnpoint(router *gin.Engine, controllCommandModel models.ControllCommandModel, connection physicaltwinconnection.Connection) {
	router.POST(controllCommandModel.ServiceUri, func(c *gin.Context) {
		var commands any
		c.BindJSON(&commands)

		response, err := connection.HandleControllCommand(controllCommandModel, commands)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, response)
	})
}
