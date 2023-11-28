package httpserver

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"os"
)

type SecureServer struct {
	server   *http.Server
	listener net.Listener
}

const SECURE_SERVER_MODE ServerMode = "secure"

func newSecureServer(url string, handler http.Handler) (*SecureServer, error) {

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

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  pool,
	}

	listener, err := net.Listen("tcp", url)
	if err != nil {
		return nil, err
	}

	// Create a new http.Server with custom configurations
	server := &http.Server{
		Addr:      listener.Addr().String(),
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	return &SecureServer{
		server:   server,
		listener: listener,
	}, nil
}

func (secureServer *SecureServer) StartServer() error {
	// Prepare the TLS certificate and private key file paths
	certFile := os.Getenv("CERT_FILE_PATH")
	keyFile := os.Getenv("KEY_FILE_PATH")

	err := secureServer.server.ServeTLS(secureServer.listener, certFile, keyFile)
	return err
}

func (secureServer *SecureServer) Addr() string {
	return secureServer.server.Addr
}
