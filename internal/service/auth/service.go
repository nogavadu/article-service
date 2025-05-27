package auth

import (
	"context"
	"fmt"
	"github.com/nogavadu/articles-service/internal/clients/auth-service/grpc"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/service"
	"log/slog"
)

type authService struct {
	log *slog.Logger

	asc *grpc.AuthServiceClient
}

func New(log *slog.Logger, asc *grpc.AuthServiceClient) service.AuthService {
	return &authService{
		log: log,
		asc: asc,
	}
}

func (s *authService) Register(ctx context.Context, authData *model.UserAuthData) (int, error) {
	const op = "authService.Register"
	log := s.log.With(slog.String("op", op))

	userId, err := s.asc.Register(ctx, converter.ToRegisterReq(authData))
	if err != nil {
		log.Error("%s: %w", op, err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (s *authService) Login(ctx context.Context, authData *model.UserAuthData) (string, error) {
	const op = "authService.Login"
	log := s.log.With(slog.String("op", op))

	token, err := s.asc.Login(ctx, converter.ToLoginReq(authData))
	if err != nil {
		log.Error("%s: %w", op, err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
