package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	authService "github.com/nogavadu/auth-service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToRegisterReq(authData *model.UserRegisterData) *authService.RegisterRequest {
	return &authService.RegisterRequest{
		Name:     wrapperspb.String(authData.Name),
		Email:    authData.Email,
		Password: authData.Password,
	}
}

func ToLoginReq(authData *model.UserAuthData) *authService.LoginRequest {
	return &authService.LoginRequest{
		Email:    authData.Email,
		Password: authData.Password,
	}
}
