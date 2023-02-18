package gateway_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type gatewayService interface {
	FindNumberPosition(ctx context.Context, number int) (int, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
}

type GatewayServer struct {
	gatewayService gatewayService
	httpListener   net.Listener
	config         *Config
}

func (s *GatewayServer) Start() error {
	err := s.newService()
	if err != nil {
		return errors.Wrap(err, "newService")
	}

	router := mux.NewRouter()
	s.GetUsersHandler().AddRoute(router)
	s.FindNumberPositionHandler().AddRoute(router)

	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		log.Errorf("Shortest URL server failed to bind to TCP port %d for listening.", s.config.Port)
		os.Exit(13)
	}

	log.Infof("Starting Shortest URL server on TCP port %d", s.config.Port)

	s.httpListener = httpListener

	err = http.Serve(httpListener, router)
	if err != nil {
		log.Errorf("Shortest URL server has shutdown: %v", err)

		return err
	}

	return nil
}

func (s *GatewayServer) newService() error {
	//TODO

	return nil
}

func (s *GatewayServer) Shutdown() {
	log.Info("Gateway server will shutdown")
	err := s.httpListener.Close()
	if err != nil {
		log.Errorf("failed while trying to close the gateway http listener: %v", err)
	}
}

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(response)
}

func New(config *Config) *GatewayServer {
	return &GatewayServer{
		config: config,
	}
}
