package grpc

import (
	"context"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	accessService "github.com/nogavadu/auth-service/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"time"
)

const (
	UserAccessLevel      = 0
	ModeratorAccessLevel = 10
	AdminAccessLevel     = 25
	CreatorAccessLevel   = 100
)

type AccessServiceClient struct {
	api accessService.AccessV1Client
	log *slog.Logger
}

func NewAccessServiceClient(
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*AccessServiceClient, error) {
	const op = "AuthServiceClient.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &AccessServiceClient{
		api: accessService.NewAccessV1Client(cc),
		log: log,
	}, nil
}

func (c *AccessServiceClient) Check(ctx context.Context, accessToken string, level int) error {
	const op = "AuthServiceClient.Check"

	md := metadata.Pairs(
		"authorization", fmt.Sprintf("Bearer %s", accessToken),
	)

	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := c.api.Check(ctx, &accessService.CheckRequest{
		RequiredLvl: uint32(level),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
