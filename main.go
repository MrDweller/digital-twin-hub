package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/MrDweller/digital-twin-hub/manufacturer"
	serviceModels "github.com/MrDweller/service-registry-connection/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	service, err := manufacturer.NewService()
	if err != nil {
		log.Panic(err)
	}
	controller := manufacturer.NewController(service)

	address := os.Getenv("ADDRESS")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Panic(err)
	}
	systemName := os.Getenv("SYSTEM_NAME")

	serviceRegistryAddress := os.Getenv("SERVICE_REGISTRY_ADDRESS")
	serviceRegistryPort, err := strconv.Atoi(os.Getenv("SERVICE_REGISTRY_PORT"))
	if err != nil {
		log.Panic(err)
	}

	digitalTwinServiceDefinition := serviceModels.ServiceDefinition{
		ServiceDefinition: "digital-twin",
		ServiceUri:        "/digital-twin",
	}
	digitalTwinCreation := manufacturer.Endpoint{
		ServiceDefinition: digitalTwinServiceDefinition,
		HttpMethod:        "POST",
		Handler:           controller.CreateDigitalTwin,
	}

	endpoints := []manufacturer.Endpoint{
		digitalTwinCreation,
	}

	manufacturer, err := manufacturer.NewManufacturer(address, port, systemName, serviceRegistryAddress, serviceRegistryPort, endpoints)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		err = manufacturer.RunManufacturerApi()
		log.Panic(err)

	}()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	log.Printf("Stopping the digital twin hub!")

	err = manufacturer.StopManufacturerApi()
	if err != nil {
		log.Panic(err)
	}
}
