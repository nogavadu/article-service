package user

import (
	"context"
	"github.com/nogavadu/articles-service/internal/clients/auth-service/grpc"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/service"
	"log/slog"
)

type userService struct {
	log *slog.Logger

	userClient *grpc.UserServiceClient
}

func New(
	log *slog.Logger,
	userClient *grpc.UserServiceClient,
) service.UserService {
	return &userService{
		log:        log,
		userClient: userClient,
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
	//TODO implement me
	panic("implement me")
}

func (s *userService) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
