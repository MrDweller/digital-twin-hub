package sensoranomalyhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	additionalservice "github.com/MrDweller/digital-twin-hub/additional-service"
	"github.com/MrDweller/digital-twin-hub/client"
	"github.com/MrDweller/digital-twin-hub/database"
	"github.com/MrDweller/digital-twin-hub/models"
	orchestratormodels "github.com/MrDweller/orchestrator-connection/models"
	"github.com/MrDweller/orchestrator-connection/orchestrator"
	servicemodels "github.com/MrDweller/service-registry-connection/models"
	serviceregistry "github.com/MrDweller/service-registry-connection/service-registry"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type SensorAnomalyHandler struct {
	router                         *gin.Engine
	notifyToHandleableAnomaliesMap map[*NotifyAnomaly]HandleableAnomaly
	anomalyHandlerService          AnomalyHandlerService
	OrchestrationConnection        orchestrator.OrchestratorConnection
	handlingSystem                 orchestratormodels.SystemDefinition
	certificateId                  string
}

func InitAnomalyHandler(anomalies []Anomaly, router *gin.Engine, SystemName string, certificateId string) ([]additionalservice.AdditionalService, error) {
	additionalServices := []additionalservice.AdditionalService{}
	handlingSystem := servicemodels.SystemDefinition{
		Address:    "localhost",
		Port:       5672,
		SystemName: SystemName,
	}

	notifyToHandleableAnomaliesMap := map[*NotifyAnomaly]HandleableAnomaly{}
	for _, anomaly := range anomalies {

		handleableAnomaly := &RabbitMQHandleableAnomaly{
			HandleableAnomalyBase: HandleableAnomalyBase{
				anomalyHandlingSystem: handlingSystem,
				Anomaly:               anomaly,
			},
			id: uuid.New(),
		}

		additionalServices = append(additionalServices, handleableAnomaly)

		notifyAnomaly := &NotifyAnomaly{
			Anomaly: anomaly,
		}
		additionalServices = append(additionalServices, notifyAnomaly)

		notifyToHandleableAnomaliesMap[notifyAnomaly] = handleableAnomaly
		fmt.Println(handleableAnomaly.GetService())
	}
	sensorAnomalyHandler, err := newSensorAnomalyHandler(
		router,
		notifyToHandleableAnomaliesMap, RabbitmqAnomalyHandlerService{
			rabbitmqAddress: handlingSystem.Address,
			rabbitmqPort:    handlingSystem.Port,
		},
		orchestratormodels.SystemDefinition{
			Address:    handlingSystem.Address,
			Port:       handlingSystem.Port,
			SystemName: handlingSystem.SystemName,
		},
		certificateId,
	)
	if err != nil {
		return nil, err
	}
	sensorAnomalyHandler.SetupEndpoints()

	return additionalServices, nil
}

func newSensorAnomalyHandler(router *gin.Engine, notifyToHandleableAnomaliesMap map[*NotifyAnomaly]HandleableAnomaly, anomalyHandlerService AnomalyHandlerService, handlingSystem orchestratormodels.SystemDefinition, certificateId string) (*SensorAnomalyHandler, error) {
	filter := bson.M{
		"certificateid": certificateId,
	}
	var certificate models.CertificateModel
	err := database.Certificate.FindOne(context.TODO(), filter).Decode(&certificate)
	if err != nil {
		return nil, err
	}

	serviceRegistryAddress := os.Getenv("SERVICE_REGISTRY_ADDRESS")
	serviceRegistryPort, err := strconv.Atoi(os.Getenv("SERVICE_REGISTRY_PORT"))
	if err != nil {
		return nil, err
	}
	serviceRegistryConnection, err := serviceregistry.NewConnection(serviceregistry.ServiceRegistry{
		Address: serviceRegistryAddress,
		Port:    serviceRegistryPort,
	}, serviceregistry.SERVICE_REGISTRY_ARROWHEAD_4_6_1, servicemodels.CertificateInfo{
		CertFilePath: certificate.CertFilePath,
		KeyFilePath:  certificate.KeyFilePath,
		Truststore:   os.Getenv("TRUSTSTORE_FILE_PATH"),
	})
	if err != nil {
		return nil, err
	}

	serviceQueryResult, err := serviceRegistryConnection.Query(servicemodels.ServiceDefinition{
		ServiceDefinition: "orchestration-service",
	})
	if err != nil {
		return nil, err
	}

	serviceQueryData := serviceQueryResult.ServiceQueryData[0]

	orchestrationConnection, err := orchestrator.NewConnection(orchestrator.Orchestrator{
		Address: serviceQueryData.Provider.Address,
		Port:    serviceQueryData.Provider.Port,
	}, orchestrator.ORCHESTRATION_ARROWHEAD_4_6_1, orchestratormodels.CertificateInfo{
		CertFilePath: certificate.CertFilePath,
		KeyFilePath:  certificate.KeyFilePath,
		Truststore:   os.Getenv("TRUSTSTORE_FILE_PATH"),
	})
	if err != nil {
		return nil, err
	}

	return &SensorAnomalyHandler{
		router:                         router,
		notifyToHandleableAnomaliesMap: notifyToHandleableAnomaliesMap,
		anomalyHandlerService:          anomalyHandlerService,
		OrchestrationConnection:        orchestrationConnection,
		handlingSystem:                 handlingSystem,
		certificateId:                  certificateId,
	}, nil

}

func (sensorAnomalyHandler *SensorAnomalyHandler) SetupEndpoints() {
	for notifyAnomaly, handleableAnomaly := range sensorAnomalyHandler.notifyToHandleableAnomaliesMap {
		sensorAnomalyHandler.addNotifyEndpoint(notifyAnomaly, handleableAnomaly)

	}
}

func (sensorAnomalyHandler *SensorAnomalyHandler) addNotifyEndpoint(notifyAnomaly *NotifyAnomaly, handleableAnomaly HandleableAnomaly) {
	sensorAnomalyHandler.router.POST(
		notifyAnomaly.GetService().ServiceDefinition.ServiceUri,
		func(c *gin.Context) {
			workDTO, err := sensorAnomalyHandler.addAnomalyAsWorkTask(handleableAnomaly)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			err = sensorAnomalyHandler.anomalyHandlerService.HandleAnomaly(
				handleableAnomaly,
				Work{
					WorkId:    workDTO.WorkId,
					ProductId: workDTO.ProductId,
				},
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.Status(http.StatusOK)
		},
	)
}

func (sensorAnomalyHandler *SensorAnomalyHandler) addAnomalyAsWorkTask(handleableAnomaly HandleableAnomaly) (*WorkDTO, error) {
	orchestrationResponse, err := sensorAnomalyHandler.OrchestrationConnection.Orchestration(
		"create-work",
		[]string{
			"HTTP-SECURE-JSON",
		},
		sensorAnomalyHandler.handlingSystem,
		orchestratormodels.AdditionalParametersArrowhead_4_6_1{
			OrchestrationFlags: map[string]bool{
				"overrideStore": true,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if len(orchestrationResponse.Response) <= 0 {
		return nil, errors.New("found no providers")
	}

	provider := orchestrationResponse.Response[0]

	createWorkDTO := CreateWorkDTO{
		EventType: handleableAnomaly.GetAnomaly().EventType,
		ProductId: sensorAnomalyHandler.handlingSystem.SystemName,
	}

	payload, err := json.Marshal(createWorkDTO)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://"+provider.Provider.Address+":"+strconv.Itoa(provider.Provider.Port)+provider.ServiceUri, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client, err := client.GetTlsClient(sensorAnomalyHandler.certificateId)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errorString := fmt.Sprintf("status: %s, body: %s", resp.Status, string(body))
		return nil, errors.New(errorString)
	}

	var workDTO WorkDTO
	err = json.Unmarshal(body, &workDTO)
	if err != nil {
		return nil, err
	}

	return &workDTO, nil
}
