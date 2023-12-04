package app

import (
	"fmt"
	"log/slog"
	"time"

	grpcapp "github.com/HironixRotifer/grpc-authorization-service/internal/app/grpc"
)

type App struct {
	GRPCsrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	//TODO: Инициализировать хранилище (storage)
	fmt.Println(storagePath)

	// TODO: Инициализировать auth service
	fmt.Println(tokenTTL.String())

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{GRPCsrv: grpcApp}
}
