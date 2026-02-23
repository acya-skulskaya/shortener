package grpc

import (
	"context"
	"errors"
	"time"

	pb "github.com/acya-skulskaya/shortener/api/shortener"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	models "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// rpc ExpandURL (URLExpandRequest) returns (URLExpandResponse);

func (s *ShortenerServer) ExpandURL(ctx context.Context, in *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	id := in.GetId()
	url, err := s.Repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrIDNotFound) {
			logger.Log.Debug("id does not exist",
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, "not found")
		} else if errors.Is(err, errorsInternal.ErrIDDeleted) {
			logger.Log.Debug("id is deleted",
				zap.Error(err),
			)
			return nil, status.Error(codes.DataLoss, "data loss")
		} else {
			logger.Log.Debug("could not get id",
				zap.Error(err),
			)
			return nil, status.Errorf(codes.Internal, "could not get url: %v", err)
		}
	}

	s.AuditPublisher.Notify(models.AuditEvent{
		Timestamp:   time.Now().Unix(),
		Action:      models.AuditEventActionTypeFollow,
		OriginalURL: url,
	})

	response := pb.URLExpandResponse_builder{
		Result: proto.String(url),
	}
	return response.Build(), nil
}
