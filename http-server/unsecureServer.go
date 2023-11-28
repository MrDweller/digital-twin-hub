package httpserver

import (
	"net"
	"net/http"
)

type UnsecureServer struct {
	server   *http.Server
	listener net.Listener
}

const UNSECURE_SERVER_MODE ServerMode = "unsecure"

func newUnSecureServer(url string, handler http.Handler) (*UnsecureServer, error) {
	// Create a listener for incoming connections
	listener, err := net.Listen("tcp", url)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: handler,
	}

	return &UnsecureServer{
		server:   server,
		listener: listener,
	}, nil
}

func (unsecureServer *UnsecureServer) StartServer() error {
	err := unsecureServer.server.Serve(unsecureServer.listener)
	return err
}

func (unsecureServer *UnsecureServer) Addr() string {
	return unsecureServer.server.Addr
}
