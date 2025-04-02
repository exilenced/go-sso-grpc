package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService authgrpc.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (app *App) Run() {
	const op = "grpcapp.Run"

	log := app.log.With(
		slog.String("op", op),
		slog.Int("port", app.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := app.gRPCServer.Serve(l); err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"

	app.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", app.port))
	app.gRPCServer.GracefulStop()
}
