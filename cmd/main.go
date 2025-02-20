package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"yandexLyceumTheme3gRPC/internal/ports"
	"yandexLyceumTheme3gRPC/internal/service"
	"yandexLyceumTheme3gRPC/pkg/api/test"
	"yandexLyceumTheme3gRPC/pkg/logger"
)

func addLogMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	ctx, _ = logger.New(ctx)
	reply, err := handler(ctx, req)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "gRPC top-level log demonstration!")
	return reply, err
}

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ordersRepo := ports.NewOrdersRepositoryInMemory()

	srv := service.New(ordersRepo)
	server := grpc.NewServer(grpc.UnaryInterceptor(addLogMiddleware))
	test.RegisterOrderServiceServer(server, srv)

	logger.GetLoggerFromCtx(ctx).Info(ctx, "listening at :50051")
	if err = server.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve", zap.Error(err))
	}
}
