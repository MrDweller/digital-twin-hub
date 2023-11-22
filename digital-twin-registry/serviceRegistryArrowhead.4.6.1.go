package digitaltwinregistry

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	digitaltwinmodels "github.com/MrDweller/digital-twin-hub/digital-twin-models"
)

const SERVICE_REGISTRY_ARROWHEAD_4_6_1 DigitalTwinRegistryImplementationType = "serviceregistry-arrowhead-4.6.1"

type ServiceRegistryArrowhead_4_6_1 struct {
	DigitalTwinRegistry
}

type RegisterServiceDTO struct {
	digitaltwinmodels.ServiceDefinition
	Interfaces     []string                           `json:"interfaces"`
	ProviderSystem digitaltwinmodels.SystemDefinition `json:"providerSystem"`
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) connect() error {

	result, err := serviceRegistry.echoServiceRegistry()
	if err != nil {
		return err
	}

	if string(result) != "Got it!" {
		return errors.New("can't establish a connection with the service registry")
	}

	return nil

}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) RegisterDigitalTwin(digitalTwinModel digitaltwinmodels.DigitalTwinModel, systemDefinition digitaltwinmodels.SystemDefinition) error {
	for _, sensedPropertyModel := range digitalTwinModel.SensedProperties {
		_, err := serviceRegistry.registerService(sensedPropertyModel.ServiceDefinition, systemDefinition)
		if err != nil {
			return err
		}

	}

	for _, controlCommandModel := range digitalTwinModel.ControlCommands {
		_, err := serviceRegistry.registerService(controlCommandModel.ServiceDefinition, systemDefinition)
		if err != nil {
			return err
		}

	}

	return nil
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) registerService(serviceDefinition digitaltwinmodels.ServiceDefinition, systemDefinition digitaltwinmodels.SystemDefinition) ([]byte, error) {
	reqisterServiceDTO := RegisterServiceDTO{
		ServiceDefinition: serviceDefinition,
		Interfaces: []string{
			"HTTP-SECURE-JSON",
		},
		ProviderSystem: systemDefinition,
	}
	payload, err := json.Marshal(reqisterServiceDTO)
	fmt.Println(string(payload))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://"+serviceRegistry.Address+":"+strconv.Itoa(serviceRegistry.Port)+"/serviceregistry/register", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client, err := serviceRegistry.getClient()
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

	if resp.StatusCode != http.StatusCreated {
		errorString := fmt.Sprintf("status: %s, body: %s", resp.Status, string(body))
		return nil, errors.New(errorString)
	}

	return body, nil
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) echoServiceRegistry() ([]byte, error) {
	req, err := http.NewRequest("GET", "https://"+serviceRegistry.Address+":"+strconv.Itoa(serviceRegistry.Port)+"/serviceregistry/echo", nil)
	if err != nil {
		return nil, err
	}

	client, err := serviceRegistry.getClient()
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

	return body, nil
}

func (serviceRegistry ServiceRegistryArrowhead_4_6_1) getClient() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(os.Getenv("CERT_FILE_PATH"), os.Getenv("KEY_FILE_PATH"))
	if err != nil {
		return nil, err
	}

	// Load truststore.p12
	truststoreData, err := os.ReadFile(os.Getenv("TRUSTSTORE_FILE_PATH"))
	if err != nil {
		return nil, err

	}

	// Extract the root certificate(s) from the truststore
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(truststoreData); !ok {
		return nil, err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				RootCAs:            pool,
				InsecureSkipVerify: false,
			},
		},
	}
	return client, nil
}
