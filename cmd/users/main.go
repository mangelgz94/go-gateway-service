package main

import (
	users_api2 "github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/users_api"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	log.Info("Starting API...")
	app := createApp()
	args := os.Args
	if len(args) > 1 && args[1] == "--version" {
		os.Exit(0)
	}
	err := app.Run(args)
	if err != nil {
		log.WithError(err).Error()
	}

	setupInterruptCloseHandler()
	log.Info("API stopping")
}

func setupInterruptCloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("API stopping")
	os.Exit(0)
}

func createApp() *cli.App {
	var config users_api2.Config
	app := cli.NewApp()
	app.Version = "1.0"
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:        "GRPC_PORT",
			EnvVars:     []string{"GRPC_PORT"},
			Value:       50051,
			Usage:       "port for grpc server",
			Destination: &config.Port,
		},
		&cli.StringFlag{
			Name:        "REPOSITORY_FILE_DIRECTORY",
			EnvVars:     []string{"REPOSITORY_FILE_DIRECTORY"},
			Value:       "../../users",
			Usage:       "file directory where file users are",
			Destination: &config.RepositoryFileDirectory,
		},
		&cli.IntFlag{
			Name:        "GRPC_SERVER_KEEPALIVE_ENFORCE_MIN_TIME",
			EnvVars:     []string{"GRPC_SERVER_KEEPALIVE_ENFORCE_MIN_TIME"},
			Value:       5,
			Usage:       "grpc server keep alive enforcement min time",
			Destination: &config.ServerKeepAliveEnforcementMinTime,
		},
		&cli.BoolFlag{
			Name:        "GRPC_SERVER_KEEPALIVE_ENFORCE_PERMIT_WITHOUT_STREAM",
			EnvVars:     []string{"GRPC_SERVER_KEEPALIVE_ENFORCE_PERMIT_WITHOUT_STREAM"},
			Usage:       "grpc server keep alive enforcement min time",
			Destination: &config.ServerKeepAlivePermitWithoutStream,
		},
		&cli.IntFlag{
			Name:        "GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_IDLE",
			EnvVars:     []string{"GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_IDLE"},
			Value:       15,
			Usage:       "grpc server keep alive max connection idle",
			Destination: &config.ServerKeepAliveMaxConnectionIdle,
		},
		&cli.IntFlag{
			Name:        "GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_AGE",
			EnvVars:     []string{"GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_AGE"},
			Value:       30,
			Usage:       "grpc server keep alive max connection age",
			Destination: &config.ServerKeepAliveMaxConnectionAge,
		},
		&cli.IntFlag{
			Name:        "GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_AGE_GRACE",
			EnvVars:     []string{"GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_AGE_GRACE"},
			Value:       5,
			Usage:       "grpc server keep alive max connection age grace",
			Destination: &config.ServerKeepAliveMaxConnectionAgeGrace,
		},
		&cli.IntFlag{
			Name:        "GRPC_SERVER_KEEP_ALIVE_TIME",
			EnvVars:     []string{"GRPC_SERVER_KEEP_ALIVE_TIME"},
			Value:       5,
			Usage:       "grpc server keep alive timeout",
			Destination: &config.ServerKeepAliveTime,
		},
		&cli.IntFlag{
			Name:        "GRPC_SERVER_KEEP_ALIVE_TIMEOUT",
			EnvVars:     []string{"GRPC_SERVER_KEEP_ALIVE_TIMEOUT"},
			Value:       1,
			Usage:       "grpc server keep alive timeout",
			Destination: &config.ServerKeepAliveTimeout,
		},
	}
	app.Before = func(ctx *cli.Context) error {
		initializeLoggingConfiguration()

		return nil
	}
	app.Action = func(ctx *cli.Context) error {

		listener, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
		if err != nil {
			return errors.Wrap(err, "net Listen")
		}
		grpcServer := users_api2.New(&config)
		grpcServer.Start()
		log.Infof("gRPC service listening on port %d", config.Port)
		if err := grpcServer.Server.Serve(listener); err != nil {
			return errors.Wrap(err, "Server Serve")
		}

		return nil
	}

	return app
}

func initializeLoggingConfiguration() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
}
