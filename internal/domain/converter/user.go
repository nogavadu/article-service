package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	authService "github.com/nogavadu/auth-service/pkg/auth_v1"
)

func ToRegisterReq(authData *model.UserAuthData) *authService.RegisterRequest {
	return &authService.RegisterRequest{
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
