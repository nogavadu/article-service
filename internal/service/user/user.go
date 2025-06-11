package user

import (
	"context"
	"errors"
	"github.com/nogavadu/articles-service/internal/clients/auth-service/grpc"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/service"
	"log/slog"
)

var (
	ErrAlreadyExists       = errors.New("article already exists")
	ErrInvalidArguments    = errors.New("invalid article arguments")
	ErrInternalServerError = errors.New("internal server error")
	ErrAccessDenied        = errors.New("access denied")
)

type userService struct {
	log *slog.Logger

	authClient   *grpc.AuthServiceClient
	accessClient *grpc.AccessServiceClient
	userClient   *grpc.UserServiceClient
}

func New(
	log *slog.Logger,
	authClient *grpc.AuthServiceClient,
	accessClient *grpc.AccessServiceClient,
	userClient *grpc.UserServiceClient,
) service.UserService {
	return &userService{
		log:          log,
		authClient:   authClient,
		accessClient: accessClient,
		userClient:   userClient,
	}
}

func (s *userService) GetById(ctx context.Context, id int) (*model.User, error) {
	const op = "userService.GetById"
	log := s.log.With(slog.String("op", op))

	user, err := s.userClient.GetById(ctx, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id int, input *model.UserUpdateInput) error {
	const op = "userService.Update"
	log := s.log.With(slog.String("op", op))

	if err := s.authClient.IsUser(ctx, id); err != nil || input.Role != nil {
		accessToken, err := s.authClient.AccessToken(ctx)
		if err != nil {
			log.Error("failed to get access token", slog.String("error", err.Error()))
			return ErrAccessDenied
		}

		err = s.accessClient.Check(ctx, accessToken, grpc.ModeratorAccessLevel)
		if err != nil {
			log.Error("failed to check access token", slog.String("error", err.Error()))
			return ErrAccessDenied
		}
	}

	err := s.userClient.Update(ctx, id, input)
	if err != nil {
		log.Error("failed to update user", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
