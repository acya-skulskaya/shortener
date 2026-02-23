package interceptors

import (
	"context"
	"fmt"
	"time"

	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	fmt.Printf("Request: method=%s\n", info.FullMethod)
	logger.Log.Info("REQUEST",
		zap.String("FullMethod", info.FullMethod),
	)

	resp, err := handler(ctx, req)

	logger.Log.Info("RESPONSE",
		zap.String("FullMethod", info.FullMethod),
		zap.Duration("duration", time.Since(start)),
		zap.Error(err),
	)

	return resp, err
}
