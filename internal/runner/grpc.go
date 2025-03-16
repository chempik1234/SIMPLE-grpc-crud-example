package runner

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"yandexLyceumTheme3gRPC/internal/ports"
	"yandexLyceumTheme3gRPC/internal/service"
	"yandexLyceumTheme3gRPC/pkg/api/test"
	"yandexLyceumTheme3gRPC/pkg/logger"
	"yandexLyceumTheme3gRPC/pkg/transport/grpc/interceptors"
)

func RunGRPC(ctx context.Context, server *grpc.Server, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "couldn't create listener on port", zap.Int("port", port), zap.Error(err))
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("listening at :%d", port))
	if err = server.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to serve grpc server", zap.Error(err))
	}
}

func CreateGRPC(ordersRepo ports.OrdersRepository) (*grpc.Server, error) {
	grpcSrv := service.New(ordersRepo)
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AddLogMiddleware))
	test.RegisterOrderServiceServer(server, grpcSrv)
	return server, nil
}
