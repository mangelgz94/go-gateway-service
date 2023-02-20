package gateway_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/models"
	protoFindNumberPositionAPI "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/proto/find-number-position-api"
	protoUsersAPI "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/proto/users-api"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/providers"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/providers/find_number_position"
	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/internal/gateway/providers/users"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	grpcUsersAPIEndpoint              = "users-api-client-connection"
	grpcFindNumberPositionAPIEndpoint = "find-number-position-api-client-connection"
)

type gatewayService interface {
	FindNumberPosition(ctx context.Context, number int) (int, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	Shutdown()
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
	grpcConfig := &providers.GrpcClientConfig{
		ClientKeepAliveTime:           s.config.GRPCClientKeepAliveTime,
		ClientKeepAliveTimeout:        s.config.GRPCClientAliveTimeout,
		ClientPermissionWithoutStream: s.config.GRPCClientPermitWithoutStream,
		ClientMaxAttempts:             s.config.GRPCClientMAxAttempts,
		ClientInitialBackoff:          s.config.GRPCClientMaxBackoff,
		ClientMaxBackoff:              s.config.GRPCClientMaxBackoff,
		ClientBackoffMultiplier:       s.config.GRPCClientBackoffMultiplier,
	}
	usersAPIGRPCClient, err := providers.NewGRPCClient(grpcConfig, s.config.GRPCUsersAddress, grpcUsersAPIEndpoint)
	if err != nil {
		log.Errorf("couldn't setup users grpc connection, error: %v", err)
		return errors.Wrap(err, "providers NewGRPCClient")
	}
	usersAPIGRPCConnection := protoUsersAPI.NewUsersAPIServiceClient(usersAPIGRPCClient.Connection)
	usersAPIProvider := users.New(usersAPIGRPCConnection, usersAPIGRPCClient)

	findNumberPositionGRPCClient, err := providers.NewGRPCClient(grpcConfig, s.config.GRPCFindNumberPositionAddress, grpcFindNumberPositionAPIEndpoint)
	if err != nil {
		log.Errorf("couldn't setup find number position grpc connection, error: %v", err)
		return errors.Wrap(err, "providers NewGRPCClient")
	}
	findNumberPositionAPIGRPCConnection := protoFindNumberPositionAPI.NewFindNumberPositionAPIServiceClient(findNumberPositionGRPCClient.Connection)
	findNumberPositionAPIProvider := find_number_position.New(findNumberPositionAPIGRPCConnection, findNumberPositionGRPCClient)

	s.gatewayService = gateway.New(usersAPIProvider, findNumberPositionAPIProvider)

	return nil
}

func (s *GatewayServer) Shutdown() {
	log.Info("Gateway server will shutdown")
	err := s.httpListener.Close()
	if err != nil {
		log.Errorf("failed while trying to close the gateway http listener: %v", err)
	}

	if s.gatewayService != nil {
		s.gatewayService.Shutdown()
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
