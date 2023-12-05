package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/HironixRotifer/grpc-authorization-service/internal/app"
	"github.com/HironixRotifer/grpc-authorization-service/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting application")

	//TODO: инициализация приложений
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	application.GRPCsrv.MustRun()

	//TODO: запустить grpc

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sn := <-stop
	log.Info("shutting down application", slog.String("signal", sn.String()))

	application.GRPCsrv.Stop()
	log.Info("application is stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
