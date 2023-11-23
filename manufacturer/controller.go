package manufacturer

import (
	"net/http"

	"github.com/MrDweller/digital-twin-hub/models"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	controller := &Controller{
		service: service,
	}

	return controller
}

func (controller *Controller) CreateDigitalTwin(c *gin.Context) {
	var digitalTwinModel models.DigitalTwinModel
	if err := c.BindJSON(&digitalTwinModel); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	systemDefinition, err := controller.service.CreateDigitalTwin(digitalTwinModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, systemDefinition)
}
