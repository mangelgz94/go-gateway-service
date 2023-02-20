package main

import (
	"os"
	"os/signal"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/gateway_api"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type server interface {
	Shutdown()
}

var gatewayServer server

func main() {
	log.Info("Starting Gateway server")

	app := createApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("Server failed, error: %v", err)
	}
	setupInterruptCloseHandler()
	log.Info("Gateway stopping")

}

func createApp() *cli.App {
	config := &gateway_api.Config{}

	app := cli.NewApp()
	app.Usage = "Run gateway server"
	app.Description = "Gateway API providers a gateway proxy to some APIs"
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:        "PORT",
			EnvVars:     []string{"PORT"},
			Value:       8090,
			Usage:       "port for API server",
			Destination: &config.Port,
		},
		&cli.IntFlag{
			Name:        "GRPC_CLIENT_KEEPALIVE_ALIVE_TIME",
			EnvVars:     []string{"GRPC_CLIENT_KEEPALIVE_ALIVE_TIME"},
			Value:       10,
			Usage:       "grpc client keep alive alive time",
			Destination: &config.GRPCClientKeepAliveTime,
		},
		&cli.IntFlag{
			Name:        "GRPC_CLIENT_KEEPALIVE_TIMEOUT",
			EnvVars:     []string{"GRPC_CLIENT_KEEPALIVE_TIMEOUT"},
			Value:       5,
			Usage:       "grpc client keep alive timeout",
			Destination: &config.GRPCClientAliveTimeout,
		},
		&cli.BoolFlag{
			Name:        "GRPC_CLIENT_PERMIT_WITHOUT_STREAM",
			EnvVars:     []string{"GRPC_CLIENT_PERMIT_WITHOUT_STREAM"},
			Usage:       "grpc client permit without stream",
			Destination: &config.GRPCClientPermitWithoutStream,
		},

		&cli.IntFlag{
			Name:        "GRPC_CLIENT_MAX_ATTEMPTS",
			EnvVars:     []string{"GRPC_CLIENT_MAX_ATTEMPTS"},
			Value:       5,
			Usage:       "grpc client max attempts",
			Destination: &config.GRPCClientMAxAttempts,
		},
		&cli.StringFlag{
			Name:        "GRPC_CLIENT_MAX_BACKOFF",
			EnvVars:     []string{"GRPC_CLIENT_MAX_BACKOFF"},
			Value:       "0.01s",
			Usage:       "grpc client max backoff",
			Destination: &config.GRPCClientMaxBackoff,
		},
		&cli.Float64Flag{
			Name:        "GRPC_CLIENT_BACKOFF_MULTIPLIER",
			EnvVars:     []string{"GRPC_CLIENT_BACKOFF_MULTIPLIER"},
			Value:       1.0,
			Usage:       "grpc client backoff multiplier",
			Destination: &config.GRPCClientBackoffMultiplier,
		},
		&cli.StringFlag{
			Name:        "GRPC_USERS_ADDRESS",
			EnvVars:     []string{"GRPC_USERS_ADDRESS"},
			Value:       "localhost:50051",
			Usage:       "grpc users address",
			Destination: &config.GRPCUsersAddress,
		},
		&cli.StringFlag{
			Name:        "GRPC_FIND_NUMBER_POSITION_ADDRESS",
			EnvVars:     []string{"GRPC_FIND_NUMBER_POSITION_ADDRESS"},
			Value:       "localhost:50052",
			Usage:       "grpc find number position address",
			Destination: &config.GRPCFindNumberPositionAddress,
		},
		&cli.StringFlag{
			Name:        "AUTH_USER",
			EnvVars:     []string{"AUTH_USER"},
			Usage:       "auth user",
			Destination: &config.AuthUser,
		},
		&cli.StringFlag{
			Name:        "AUTH_PASSWORD",
			EnvVars:     []string{"AUTH_PASSWORD"},
			Usage:       "auth password",
			Destination: &config.AuthPassword,
		},
	}

	app.Action = func(context *cli.Context) error {
		apiServer := gateway_api.New(config)

		gatewayServer = apiServer
		return apiServer.Start()
	}
	app.ExitErrHandler = func(cCtx *cli.Context, err error) {
		shutdown()
	}

	return app

}

func setupInterruptCloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("URL Shortener server stopping")
	shutdown()
	os.Exit(0)
}

func shutdown() {
	if gatewayServer != nil {
		gatewayServer.Shutdown()
	}
}
