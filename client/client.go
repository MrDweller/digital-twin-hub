package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"

	"github.com/MrDweller/digital-twin-hub/database"
	"github.com/MrDweller/digital-twin-hub/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTlsClient(certificateId string) (*http.Client, error) {
	filter := bson.M{
		"certificateid": certificateId,
	}
	var certificate models.CertificateModel
	err := database.Certificate.FindOne(context.TODO(), filter).Decode(&certificate)
	if err != nil {
		return nil, err
	}
	cert, err := tls.LoadX509KeyPair(certificate.CertFilePath, certificate.KeyFilePath)
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
