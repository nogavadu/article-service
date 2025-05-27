package grpc

import (
	"context"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	authService "github.com/nogavadu/auth-service/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type AuthServiceClient struct {
	api authService.AuthV1Client
	log *slog.Logger
}

func NewAuthServiceClient(
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*AuthServiceClient, error) {
	const op = "AccessServiceClient.New"

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

	return &AuthServiceClient{
		api: authService.NewAuthV1Client(cc),
		log: log,
	}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields)
	})
}

func (c *AuthServiceClient) Register(ctx context.Context, request *authService.RegisterRequest) (int, error) {
	const op = "AuthServiceClient.Register"

	resp, err := c.api.Register(ctx, request)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(resp.UserId), nil
}

func (c *AuthServiceClient) Login(ctx context.Context, request *authService.LoginRequest) (string, error) {
	const op = "AuthServiceClient.Login"

	resp, err := c.api.Login(ctx, request)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.RefreshToken, nil
}

func (c *AuthServiceClient) RefreshToken(ctx context.Context) (string, error) {
	const op = "AuthServiceClient.GetRefreshToken"

	resp, err := c.api.GetRefreshToken(ctx, &authService.GetRefreshTokenRequest{
		RefreshToken: ctx.Value("authorization").(string),
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.RefreshToken, nil
}

func (c *AuthServiceClient) AccessToken(ctx context.Context) (string, error) {
	const op = "AuthServiceClient.AccessToken"

	c.log.Info(ctx.Value("authorization").(string))
	resp, err := c.api.GetAccessToken(ctx, &authService.GetAccessTokenRequest{
		RefreshToken: ctx.Value("authorization").(string),
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.AccessToken, nil
}
