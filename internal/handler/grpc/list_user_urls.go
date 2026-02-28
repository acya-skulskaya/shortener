package grpc

import (
	"context"

	pb "github.com/acya-skulskaya/shortener/api/shortener"
	"github.com/acya-skulskaya/shortener/internal/logger"
	authService "github.com/acya-skulskaya/shortener/internal/service/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ShortenerServer) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*pb.UserURLsResponse, error) {
	userID, ok := ctx.Value(authService.AuthContextKey(authService.AuthContextKeyUserID)).(string)
	if !ok {
		logger.Log.Debug("could not get userID from context")
		return nil, status.Error(codes.Unauthenticated, "could not get userID from context")
	}

	list, err := s.Repo.GetUserUrls(ctx, userID)
	if err != nil {
		logger.Log.Debug("error getting user urls", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "could not get user urls: %v", err)
	}

	if len(list) == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	urlDataList := make([]*pb.URLData, 0, len(list))
	for _, item := range list {
		data := pb.URLData_builder{
			ShortUrl:    proto.String(item.ShortURL),
			OriginalUrl: proto.String(item.OriginalURL),
		}.Build()
		urlDataList = append(urlDataList, data)
	}

	response := pb.UserURLsResponse_builder{
		Url: urlDataList,
	}
	return response.Build(), nil
}
