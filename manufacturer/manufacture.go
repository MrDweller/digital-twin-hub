package manufacturer

import (
	"fmt"
	"log"
	"os"

	"github.com/MrDweller/service-registry-connection/models"
	serviceregistry "github.com/MrDweller/service-registry-connection/service-registry"
	"github.com/gin-gonic/gin"
)

type Manufacturer struct {
	models.SystemDefinition
	ServiceRegistryConnection serviceregistry.ServiceRegistryConnection
	Endpoints                 []Endpoint
}

type Endpoint struct {
	models.ServiceDefinition
	HttpMethod string
	Handler    gin.HandlerFunc
}

func NewManufacturer(address string, port int, systemName string, serviceRegistryAddress string, serviceRegistryPort int, endpoints []Endpoint) (*Manufacturer, error) {
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
		Endpoints:                 endpoints,
	}, nil
}

func (manufacturer Manufacturer) RunManufacturerApi() error {
	router := gin.Default()

	manufacturer.setupEnpoints(router)

	url := fmt.Sprintf("%s:%d", manufacturer.Address, manufacturer.Port)
	log.Printf("Starting digital twin framework on: http://%s", url)

	err := router.Run(url)
	return err

}

func (manufacturer Manufacturer) setupEnpoints(router *gin.Engine) {
	for _, endpoint := range manufacturer.Endpoints {
		router.Handle(endpoint.HttpMethod, endpoint.ServiceUri, endpoint.Handler)

		manufacturer.ServiceRegistryConnection.RegisterService(endpoint.ServiceDefinition, manufacturer.SystemDefinition)

	}

}

func (manufacturer Manufacturer) StopManufacturerApi() error {
	log.Printf("Unregistering the manufacturer services from the service registry!")
	for _, endpoint := range manufacturer.Endpoints {
		err := manufacturer.ServiceRegistryConnection.UnRegisterService(endpoint.ServiceDefinition, manufacturer.SystemDefinition)
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
