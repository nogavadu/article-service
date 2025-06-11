package grpc

import (
	"context"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	userService "github.com/nogavadu/auth-service/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log/slog"
	"time"
)

type UserServiceClient struct {
	api userService.UserV1Client
	log *slog.Logger
}

func NewUserServiceClient(
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*UserServiceClient, error) {
	const op = "AuthService.NewUserServiceClient"

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

	return &UserServiceClient{
		api: userService.NewUserV1Client(cc),
		log: log,
	}, nil
}

func (c *UserServiceClient) GetById(ctx context.Context, userId int) (*model.User, error) {
	res, err := c.api.GetById(ctx, &userService.GetByIdRequest{Id: int64(userId)})
	if err != nil {
		return nil, err
	}
	user := res.GetUser()
	userInfo := user.GetInfo()

	return &model.User{
		Id: int(user.GetId()),
		UserInfo: model.UserInfo{
			Name:   converter.ProtoStringToPtrString(userInfo.GetName()),
			Email:  userInfo.GetEmail(),
			Avatar: converter.ProtoStringToPtrString(userInfo.GetAvatar()),
			Role:   userInfo.GetRole(),
		},
	}, nil
}

func (c *UserServiceClient) Update(ctx context.Context, userId int, updateInput *model.UserUpdateInput) error {
	_, err := c.api.Update(ctx, &userService.UpdateRequest{
		Id: int64(userId),
		UpdateInput: &userService.UserUpdateInput{
			Name:   StrPtrToProtoString(updateInput.Name),
			Email:  StrPtrToProtoString(updateInput.Email),
			Avatar: StrPtrToProtoString(updateInput.Avatar),
			Role:   StrPtrToProtoString(updateInput.Role),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func StrPtrToProtoString(ptr *string) *wrapperspb.StringValue {
	if ptr == nil {
		return nil
	}

	return wrapperspb.String(*ptr)
}

func (c *UserServiceClient) Delete(ctx context.Context, userId int) error {
	panic("Implement me")

}
