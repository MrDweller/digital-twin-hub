package manufacturer

import (
	"net/http"

	_ "github.com/MrDweller/digital-twin-hub/docs"
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

// Create a new digital twin
// @Summary      Create a new digital twin
// @Description  Create a new digital twin based on the given JSON object. This will create a connection to the physical twin based on the connection info given, this will also include generating endpoints to controll and view sensed data.
// @Tags         Management
// @Produce      json
// @Param        DigitalTwinModel  body       DigitalTwinModelDTO  true  "DigitalTwinModel JSON"
// @Success      200 {object} SystemDefinitionDTO
// @Router       /digital-twin [post]
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
