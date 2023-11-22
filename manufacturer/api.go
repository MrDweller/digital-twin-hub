package manufacturer

import (
	"github.com/gin-gonic/gin"
)

func RunManufacturerApi(url string) error {
	router := gin.Default()

	service, err := NewService()
	if err != nil {
		return err
	}
	controller := NewController(service)

	router.POST("/digital-twin", controller.CreateDigitalTwin)

	err = router.Run(url)
	return err
}
