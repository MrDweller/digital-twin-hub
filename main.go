package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/MrDweller/digital-twin-hub/manufacturer"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	url := fmt.Sprintf("%s:%d", os.Getenv("ADDRESS"), port)
	log.Printf("Starting digital twin framework on: http://%s", url)

	err = manufacturer.RunManufacturerApi(url)
	log.Panic(err)
}
