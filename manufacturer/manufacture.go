package manufacturer

import (
	"fmt"
	"log"
	"os"

	_ "github.com/MrDweller/digital-twin-hub/docs"
	httpserver "github.com/MrDweller/digital-twin-hub/http-server"
	"github.com/MrDweller/service-registry-connection/models"
	serviceregistry "github.com/MrDweller/service-registry-connection/service-registry"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Manufacturer struct {
	models.SystemDefinition
	ServiceRegistryConnection serviceregistry.ServiceRegistryConnection
	Services                  []models.ServiceDefinition
	Service                   *Service
}

func NewManufacturer(address string, port int, systemName string, serviceRegistryAddress string, serviceRegistryPort int, services []models.ServiceDefinition) (*Manufacturer, error) {
	system := models.SystemDefinition{
		Address:            address,
		Port:               port,
		SystemName:         systemName,
		AuthenticationInfo: os.Getenv("AUTHENTICATION_INFO"),
	}

	serviceRegistryConnection, err := serviceregistry.NewConnection(serviceregistry.ServiceRegistry{
		Address: serviceRegistryAddress,
		Port:    serviceRegistryPort,
	}, serviceregistry.SERVICE_REGISTRY_ARROWHEAD_4_6_1, models.CertificateInfo{
		CertFilePath: os.Getenv("CERT_FILE_PATH"),
		KeyFilePath:  os.Getenv("KEY_FILE_PATH"),
		Truststore:   os.Getenv("TRUSTSTORE_FILE_PATH"),
	})
	if err != nil {
		return nil, err
	}

	service, err := NewService()
	if err != nil {
		log.Panic(err)
	}

	return &Manufacturer{
		SystemDefinition:          system,
		ServiceRegistryConnection: serviceRegistryConnection,
		Services:                  services,
		Service:                   service,
	}, nil
}

func (manufacturer Manufacturer) RunManufacturer() error {
	err := manufacturer.Service.startAllSavedDigitalTwins()
	if err != nil {
		return err
	}

	router := gin.Default()

	url := fmt.Sprintf("%s:%d", manufacturer.Address, manufacturer.Port)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	manufacturer.setupEnpoints(router, url)
	manufacturer.registerServices()

	log.Printf("Starting digital twin framework on: https://%s", url)
	log.Printf("Swagger documentation is available on: https://%s", url+"/docs/index.html")

	server, err := httpserver.NewServer(url, router)
	if err != nil {
		return err
	}
	err = server.StartServer()
	return err

}

func (manufacturer Manufacturer) setupEnpoints(router *gin.Engine, url string) error {
	controller := NewController(manufacturer.Service)

	router.POST("/digital-twin", AdminAuthorization, controller.CreateDigitalTwin)
	router.DELETE("/digital-twin", AdminAuthorization, controller.DeleteDigitalTwin)

	return nil
}

func (manufacturer Manufacturer) registerServices() {
	for _, service := range manufacturer.Services {
		manufacturer.ServiceRegistryConnection.RegisterService(service, manufacturer.SystemDefinition)
	}

}

func (manufacturer Manufacturer) StopManufacturer() error {
	log.Printf("Unregistering the manufacturer services from the service registry!")
	for _, service := range manufacturer.Services {
		err := manufacturer.ServiceRegistryConnection.UnRegisterService(service, manufacturer.SystemDefinition)
		if err != nil {
			return err
		}
	}

	err := manufacturer.ServiceRegistryConnection.UnRegisterSystem(manufacturer.SystemDefinition)
	if err != nil {
		return err
	}

	err = manufacturer.Service.stopAllSavedDigitalTwins()
	if err != nil {
		return err
	}

	return nil

}
