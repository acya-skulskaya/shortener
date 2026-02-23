package interceptors

import (
	"context"
	"errors"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	authService "github.com/acya-skulskaya/shortener/internal/service/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var token string
	var err error
	if ok && len(md[authService.AuthMetadataKeyName]) > 0 {
		token = md[authService.AuthMetadataKeyName][0]
	} else {
		token, err = authService.BuildJWTString()
		if err != nil {
			logger.Log.Debug("AuthUnaryInterceptor: could not create token string", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "could not create token string: %v", err)
		}
		if err = grpc.SetHeader(ctx, metadata.Pairs(authService.AuthMetadataKeyName, token)); err != nil {
			return nil, status.Errorf(codes.Internal, "could not set gRPC header: %v", err)
		}
	}

	userID, err := authService.GetUserID(token)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrTokenIsNotValid) {
			return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: could not validate token: %v", err)
		} else {
			logger.Log.Debug("error getting user id from auth token", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "could not set user grom token: %v", err)
		}
	}

	if userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: user is empty")
	}

	logger.Log.Info("AuthUnaryInterceptor: got user id", zap.String("userID", userID))
	ctx = context.WithValue(ctx, authService.AuthContextKey(authService.AuthContextKeyUserID), userID)

	return handler(ctx, req)
}
