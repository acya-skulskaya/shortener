package grpc

import (
	"context"
	"errors"
	"time"

	pb "github.com/acya-skulskaya/shortener/api/shortener"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	models "github.com/acya-skulskaya/shortener/internal/model/json"
	authService "github.com/acya-skulskaya/shortener/internal/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (s *ShortenerServer) ShortenURL(ctx context.Context, in *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	url := in.GetUrl()
	userID, ok := ctx.Value(authService.AuthContextKey(authService.AuthContextKeyUserID)).(string)
	if !ok {
		logger.Log.Debug("could not get userID from context")
		return nil, status.Error(codes.Unauthenticated, "could not get userID from context")
	}

	id, err := s.Repo.Store(ctx, url, userID)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrConflictOriginalURL) && len(id) > 0 {
			return nil, status.Errorf(codes.AlreadyExists, "url already shortened, id %s", id)
		}

		return nil, status.Errorf(codes.Internal, "could not create short url: %v", err)
	}

	s.AuditPublisher.Notify(models.AuditEvent{
		Timestamp:   time.Now().Unix(),
		Action:      models.AuditEventActionTypeShorten,
		UserID:      userID,
		OriginalURL: url,
	})

	response := pb.URLShortenResponse_builder{
		Result: proto.String(id),
	}
	return response.Build(), nil
}
