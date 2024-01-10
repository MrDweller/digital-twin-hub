package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database
var DigitalTwin *mongo.Collection
var SensorData *mongo.Collection

func InitDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	dbConnectionString := fmt.Sprintf(os.Getenv("MONGO_DB_CONNECTION_STRING"))

	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnectionString))
	if err != nil {
		return err
	}
	Database = Client.Database("DigitalTwinHub")
	DigitalTwin = Database.Collection("DigitalTwin")
	SensorData = Database.Collection("SensorData")

	return nil
}
