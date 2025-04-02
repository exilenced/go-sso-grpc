package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	tokenTTL time.Duration,
	storage *postgres.Storage,
) *App {
	storage, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, tokenTTL)
	//TODO
	grpcApp := grpcapp.New(log, authService, cfg.GRPC.Port)
	return &App{
		GRPCSrv: grpcApp,
	}
}
