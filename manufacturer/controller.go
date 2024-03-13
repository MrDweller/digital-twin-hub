package manufacturer

import (
	"log"
	"net/http"
	"strconv"

	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	_ "github.com/MrDweller/digital-twin-hub/docs"
	sensoranomalyhandler "github.com/MrDweller/digital-twin-hub/sensor-anomaly-handler"
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
// @Param        DigitalTwin  body       DigitalTwinDTO  true  "DigitalTwinDTO JSON"
// @Success      200 {object} SystemDefinitionDTO
// @Router       /create-digital-twin [post]
func (controller *Controller) CreateDigitalTwin(c *gin.Context) {

	var digitalTwinDto DigitalTwinDTO
	if err := c.BindJSON(&digitalTwinDto); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	additionalServices := []additionalservice.AdditionalService{}
	router := gin.New()

	anomalyServices, err := sensoranomalyhandler.InitAnomalyHandler(
		mapAnomaliesDtoToAnomalies(digitalTwinDto.HandleableAnomalies),
		router,
		digitalTwinDto.SystemName,
		digitalTwinDto.CertificateId,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	additionalServices = append(additionalServices, anomalyServices...)

	systemDefinition, err := controller.service.CreateDigitalTwin(mapDigitalTwinDtoToDigitalTwinModel(digitalTwinDto, additionalServices), uuid.New(), digitalTwinDto.SystemName, router)
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

// Upload certificate files
// @Summary      Upload certificate files as zip
// @Description  Upload certificate files as zip, to be used by a digital twin. Takes cert.pem and key.pem files and gives `certId`.
// @Tags         Management
// @Produce      json
// @Param        cert formData  file  true  "Cert file"
// @Param        key formData  file  true  "Key file"
// @Success      200 {object} CertificateDTO
// @Router       /upload-certificates [post]
func (controller *Controller) UploadCertificates(c *gin.Context) {

	certFile, err := c.FormFile("cert")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	keyFile, err := c.FormFile("key")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	CertificateModel := controller.service.UploadCertificates(certFile.Filename, keyFile.Filename)

	err = c.SaveUploadedFile(certFile, CertificateModel.CertFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Printf("saved cert file: %s\n", CertificateModel.CertFilePath)

	err = c.SaveUploadedFile(keyFile, CertificateModel.KeyFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Printf("saved key file: %s", CertificateModel.CertFilePath)

	c.JSON(http.StatusOK, CertificateDTO{
		CertificateId: CertificateModel.CertificateId,
	})
}
