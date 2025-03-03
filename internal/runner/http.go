package runner

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"yandexLyceumTheme3gRPC/pkg/api/test"
	"yandexLyceumTheme3gRPC/pkg/logger"
)

func CreateHTTP(ctx context.Context, grpcEndpoint string, port int) (*http.Server, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := test.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		mux,
		grpcEndpoint,
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register gRPC Gateway: %w", err)
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}, nil
}

func RunHTTP(ctx context.Context, srv *http.Server) {
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("listening at %s", srv.Addr))
	if err := srv.ListenAndServe(); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve gateway", zap.Error(err))
	}
}
