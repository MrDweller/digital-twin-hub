package manufacturer

import (
	"fmt"
	"log"
	"os"

	httpserver "github.com/MrDweller/digital-twin-hub/http-server"
	"github.com/MrDweller/service-registry-connection/models"
	serviceregistry "github.com/MrDweller/service-registry-connection/service-registry"
	"github.com/gin-gonic/gin"
)

type Manufacturer struct {
	models.SystemDefinition
	ServiceRegistryConnection serviceregistry.ServiceRegistryConnection
	Services                  []models.ServiceDefinition
}

func NewManufacturer(address string, port int, systemName string, serviceRegistryAddress string, serviceRegistryPort int, services []models.ServiceDefinition) (*Manufacturer, error) {
	system := models.SystemDefinition{
		Address:    address,
		Port:       port,
		SystemName: systemName,
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

	return &Manufacturer{
		SystemDefinition:          system,
		ServiceRegistryConnection: serviceRegistryConnection,
		Services:                  services,
	}, nil
}

func (manufacturer Manufacturer) RunManufacturerApi() error {
	router := gin.Default()

	manufacturer.setupEnpoints(router)
	manufacturer.registerServices()

	url := fmt.Sprintf("%s:%d", manufacturer.Address, manufacturer.Port)
	log.Printf("Starting digital twin framework on: http://%s", url)

	server, err := httpserver.NewServer(url, router)
	if err != nil {
		return err
	}
	err = server.StartServer()
	return err

}

func (manufacturer Manufacturer) setupEnpoints(router *gin.Engine) error {
	service, err := NewService()
	if err != nil {
		log.Panic(err)
	}
	controller := NewController(service)

	router.POST("/digital-twin", AdminAuthorization, controller.CreateDigitalTwin)

	return nil
}

func (manufacturer Manufacturer) registerServices() {
	for _, service := range manufacturer.Services {
		manufacturer.ServiceRegistryConnection.RegisterService(service, manufacturer.SystemDefinition)
	}

}

func (manufacturer Manufacturer) StopManufacturerApi() error {
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

	return nil

}
