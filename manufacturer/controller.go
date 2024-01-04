package manufacturer

import (
	"net/http"
	"strconv"

	_ "github.com/MrDweller/digital-twin-hub/docs"
	"github.com/MrDweller/digital-twin-hub/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Router       /create-digital-twin [post]
func (controller *Controller) CreateDigitalTwin(c *gin.Context) {

	var digitalTwinModel models.DigitalTwinModel
	if err := c.BindJSON(&digitalTwinModel); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	systemDefinition, err := controller.service.CreateDigitalTwin(digitalTwinModel, uuid.New())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, systemDefinition)
}

// Delete a digital twin
// @Summary      Delete a digital twin
// @Description  Delete a digital twin based on the given address and port.
// @Tags         Management
// @Param        address 				query       string  true  "address"
// @Param        port   				query       string  true  "port "
// @Success      200
// @Router       /remove-digital-twin [delete]
func (controller *Controller) DeleteDigitalTwin(c *gin.Context) {
	address := c.Query("address")
	port, err := strconv.Atoi(c.Query("port"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = controller.service.DeleteDigitalTwin(address, port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
