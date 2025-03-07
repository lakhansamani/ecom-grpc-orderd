package service

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/sdk/trace"
	otrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"

	order "github.com/lakhansamani/ecom-grpc-apis/order/v1"
	user "github.com/lakhansamani/ecom-grpc-apis/user/v1"

	"github.com/lakhansamani/ecom-grpc-orderd/db"
)

type Config struct {
	UserServiceAddress string
}

type Dependencies struct {
	// Add dependencies here
	DBProvider db.Provider
	// UserService user.Service
	UserService   user.UserServiceClient
	TraceProvider *trace.TracerProvider
}

// Service implements the Order service.
type Service interface {
	order.OrderServiceServer
}

type service struct {
	Config
	Dependencies

	trace otrace.Tracer
}

// New creates a new Order service.
func New(cfg Config, deps Dependencies) Service {
	trace := deps.TraceProvider.Tracer("service")
	return &service{
		Config:       cfg,
		Dependencies: deps,
		trace:        trace,
	}
}

// authorize verifies user using the user service and gets userID
func (s *service) authorize(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}

	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 {
		return "", errors.New("missing authorization token")
	}
	token := authHeader[0]

	// add token to outgoing context
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	// Call user service to get user
	userResp, err := s.UserService.Me(ctx, &user.MeRequest{})
	if err != nil {
		return "", err
	}
	return userResp.GetUser().GetId(), nil
}
