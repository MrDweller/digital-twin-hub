package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/MrDweller/digital-twin-hub/database"
	_ "github.com/MrDweller/digital-twin-hub/docs"
	"github.com/MrDweller/digital-twin-hub/manufacturer"
	serviceModels "github.com/MrDweller/service-registry-connection/models"
	"github.com/joho/godotenv"
)

// @title          Digital Twin Hub
// @version        1.0
// @description    This page shows the REST interfaces offered by the Digital Twin Hub.
// @contact.url    https://github.com/MrDweller/digital-twin-hub

// @host      localhost:8080
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	err = database.InitDatabase()
	if err != nil {
		log.Panic(err)
	}

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

	manufacturer, err := manufacturer.NewManufacturer(address, port, systemName, serviceRegistryAddress, serviceRegistryPort, []serviceModels.ServiceDefinition{
		digitalTwinServiceDefinition,
	})
	if err != nil {
		log.Panic(err)
	}

	go func() {
		err = manufacturer.RunManufacturer()
		log.Panic(err)

	}()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	log.Printf("Stopping the digital twin hub!")

	err = manufacturer.StopManufacturer()
	if err != nil {
		log.Panic(err)
	}
}
