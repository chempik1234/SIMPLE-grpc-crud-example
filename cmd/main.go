package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"yandexLyceumTheme3gRPC/internal/config"
	"yandexLyceumTheme3gRPC/internal/ports/adapters"
	"yandexLyceumTheme3gRPC/internal/runner"
	"yandexLyceumTheme3gRPC/pkg/logger"
	"yandexLyceumTheme3gRPC/pkg/postgres"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, _ = logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to load config", zap.Error(err))
	}

	pgCfg := cfg.Postgres

	pool, err := postgres.New(ctx, pgCfg)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to database", zap.Error(err))
	}

	ordersRepo := adapters.NewOrdersRepositoryPostgres(pool)

	grpcServer, err := runner.CreateGRPC(ordersRepo)
	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}
	httpServer, err := runner.CreateHTTP(ctx, fmt.Sprintf("localhost:%d", cfg.GRPCPort), cfg.HTTPPort)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "could not register http->grpc gateway handler", zap.Error(err))
	}

	go runner.RunGRPC(ctx, grpcServer, cfg.GRPCPort)
	// go runner.RunHTTP(ctx, httpServer)

	/*
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
	*/

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		httpServer.Shutdown(ctx)
		pool.Close()
		logger.GetLoggerFromCtx(ctx).Info(ctx, "server stopped")
	}

	/*
		cancelCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		grpcServer.GracefulStop()
		httpServer.Shutdown(cancelCtx)
		log.Println("Server Stopped")
	*/
}
