package manufacturer

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/MrDweller/digital-twin-hub/database"
	digitaltwin "github.com/MrDweller/digital-twin-hub/digital-twin"
	digitaltwinregistry "github.com/MrDweller/digital-twin-hub/digital-twin-registry"
	"github.com/MrDweller/digital-twin-hub/models"
	serviceModels "github.com/MrDweller/service-registry-connection/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	digitalTwinRegistryAddress string
	digitalTwinRegistryPort    int
}

func NewService() (*Service, error) {
	digitalTwinRegistryPort, err := strconv.Atoi(os.Getenv("DIGITAL_TWIN_REGISTRY_PORT"))
	if err != nil {
		return nil, err
	}

	service := &Service{
		digitalTwinRegistryAddress: os.Getenv("DIGITAL_TWIN_REGISTRY_ADDRESS"),
		digitalTwinRegistryPort:    digitalTwinRegistryPort,
	}
	return service, nil
}

func (service Service) CreateDigitalTwin(digitalTwinModel models.DigitalTwinModel, digitalTwinId uuid.UUID, systemName string, router *gin.Engine) (*serviceModels.SystemDefinition, error) {

	filter := bson.M{
		"certificateid": digitalTwinModel.CertificateId,
	}
	var certificate models.CertificateModel
	err := database.Certificate.FindOne(context.TODO(), filter).Decode(&certificate)
	if err != nil {
		return nil, err
	}

	certFilePath := certificate.CertFilePath
	keyFilePath := certificate.KeyFilePath

	digitalTwin, err := digitaltwin.NewDigitalTwin(
		digitalTwinModel,
		digitaltwinregistry.DigitalTwinRegistry{
			Address: service.digitalTwinRegistryAddress,
			Port:    service.digitalTwinRegistryPort,
			CertificateInfo: serviceModels.CertificateInfo{
				CertFilePath: certFilePath,
				KeyFilePath:  keyFilePath,
				Truststore:   os.Getenv("TRUSTSTORE_FILE_PATH"),
			},
		},
		digitalTwinId,
		systemName,
		router,
	)
	if err != nil {
		return nil, err
	}

	// database.DigitalTwin.InsertOne(context.Background(), digitalTwin)

	systemDefinition, err := digitalTwin.StartDigitalTwin()
	if err != nil {
		return nil, err
	}

	return systemDefinition, nil
}

func (service Service) DeleteDigitalTwin(address string, port int) error {
	filter := bson.M{
		"systemdefinition.address": address,
		"systemdefinition.port":    port,
	}
	digitalTwins, err := service.getSavedDigitalTwins(filter)
	database.DigitalTwin.DeleteMany(context.Background(), filter)

	if err != nil {
		return err
	}

	for _, digitalTwin := range digitalTwins {
		database.SensorData.DeleteMany(context.Background(), bson.M{
			"digitaltwinid": digitalTwin.DigitalTwinId,
		})
		// err := service.digitalTwinRegistryConnection.UnRegisterDigitalTwin(didigitalTwin.DigitalTwinModel, didigitalTwin.SystemDefinition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service Service) startAllSavedDigitalTwins() error {
	digitalTwins, err := service.getSavedDigitalTwins(bson.M{})
	if err != nil {
		return err
	}

	for _, digitalTwin := range digitalTwins {
		filter := bson.M{
			"systemdefinition.address": digitalTwin.SystemDefinition.Address,
			"systemdefinition.port":    digitalTwin.SystemDefinition.Port,
		}
		_, err := database.DigitalTwin.DeleteMany(context.Background(), filter)
		if err != nil {
			return err
		}
	}

	for _, digitalTwin := range digitalTwins {
		router := gin.New()
		_, err = service.CreateDigitalTwin(digitalTwin.DigitalTwinModel, digitalTwin.DigitalTwinId, digitalTwin.SystemDefinition.SystemName, router)
		if err != nil {
			return err
		}
	}
	return nil

}

func (service Service) stopAllSavedDigitalTwins() error {
	digitalTwins, err := service.getSavedDigitalTwins(bson.M{})
	if err != nil {
		return err
	}

	for _, digitalTwin := range digitalTwins {
		log.Println(digitalTwin)
		// err := service.digitalTwinRegistryConnection.UnRegisterDigitalTwin(digitalTwin.DigitalTwinModel, digitalTwin.SystemDefinition)
		if err != nil {
			return err
		}
	}
	return nil

}

func (service Service) getSavedDigitalTwins(filter interface{}, opts ...*options.FindOptions) ([]digitaltwin.DigitalTwin, error) {
	result, err := database.DigitalTwin.Find(context.Background(), filter, opts...)
	if err != nil {
		return nil, err
	}

	var digitalTwins []digitaltwin.DigitalTwin
	if err := result.All(context.Background(), &digitalTwins); err != nil {
		return nil, err
	}

	return digitalTwins, nil
}

func (service Service) UploadCertificates(certFileName string, keyFileName string) models.CertificateModel {
	certificateId := uuid.New()

	certificateDirectoryPath := fmt.Sprintf("certificates/digital-twins/%s", certificateId)

	certFilePath := fmt.Sprintf("%s/%s", certificateDirectoryPath, certFileName)

	keyFilePath := fmt.Sprintf("%s/%s", certificateDirectoryPath, keyFileName)

	certificateModel := models.CertificateModel{
		CertificateId: certificateId,
		CertFilePath:  certFilePath,
		KeyFilePath:   keyFilePath,
	}

	database.Certificate.InsertOne(context.Background(), certificateModel)

	return certificateModel

}
