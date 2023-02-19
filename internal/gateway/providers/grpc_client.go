package providers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

type GrpcClientConfig struct {
	ClientKeepAliveTime           int
	ClientKeepAliveTimeout        int
	ClientPermissionWithoutStream bool
	ClientMaxAttempts             int
	ClientInitialBackoff          string
	ClientMaxBackoff              string
	ClientBackoffMultiplier       float64
}

type GrpcClient struct {
	Connection   *grpc.ClientConn
	endPointName string
}

func (c *GrpcClient) CloseConnection() {
	err := c.Connection.Close()
	if err != nil {
		log.Errorf("%s: error closing response stream, %v", c.endPointName, err)
	}
	log.Infof("Client connection %s has been closed", c.endPointName)
}

func NewGRPCClient(config *GrpcClientConfig, address string, endPointName string) (*GrpcClient, error) {
	retryPolicyTemplate := `{
						  "methodConfig": [
							{
							  "name": [{"service":"grpc.users"}],
							  "waitForReady": true,
							  "retryPolicy": {
								"MaxAttempts": %v,
								"InitialBackoff": "%v",
								"MaxBackoff": "%v",
								"BackoffMultiplier": %v,
								"RetryableStatusCodes": ["UNAVAILABLE"]
							  }
							}
						  ]
						}`
	retryPolicy := fmt.Sprintf(retryPolicyTemplate, config.ClientMaxAttempts, config.ClientInitialBackoff, config.ClientMaxBackoff, config.ClientBackoffMultiplier)

	dialConfig := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(config.ClientKeepAliveTime) * time.Second,
			Timeout:             time.Duration(config.ClientKeepAliveTimeout) * time.Second,
			PermitWithoutStream: config.ClientPermissionWithoutStream,
		}),
		grpc.WithDefaultServiceConfig(retryPolicy),
	}

	connection, err := grpc.Dial(address, dialConfig...)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup %s grpc connection", address)
	}

	return &GrpcClient{Connection: connection, endPointName: endPointName}, nil
}
