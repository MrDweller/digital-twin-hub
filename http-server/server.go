package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Server interface {
	StartServer() error
	Addr() string
}

type ServerMode string

func NewServer(url string, handler http.Handler) (Server, error) {
	serverMode := os.Getenv("SERVER_MODE")

	switch ServerMode(serverMode) {
	case UNSECURE_SERVER_MODE:
		unsecureServer, err := newUnSecureServer(url, handler)
		if err != nil {
			return nil, err
		}
		return unsecureServer, nil
	case SECURE_SERVER_MODE:
		secureServer, err := newSecureServer(url, handler)
		if err != nil {
			return nil, err
		}
		return secureServer, nil
	default:
		errorString := fmt.Sprintf("the server mode %s, has no implementation", serverMode)
		return nil, errors.New(errorString)
	}
}
