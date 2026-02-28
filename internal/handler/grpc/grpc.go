package grpc

import (
	pb "github.com/acya-skulskaya/shortener/api/shortener"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	interfaces "github.com/acya-skulskaya/shortener/internal/repository/interface"
)

type ShortenerServer struct {
	pb.UnimplementedShortenerServiceServer

	Repo           interfaces.ShortURLRepository
	AuditPublisher publisher.Publisher
}
