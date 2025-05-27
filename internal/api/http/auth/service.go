package auth

import "github.com/nogavadu/articles-service/internal/service"

type Implementation struct {
	authServ service.AuthService
}

func New(authServ service.AuthService) *Implementation {
	return &Implementation{
		authServ: authServ,
	}
}
