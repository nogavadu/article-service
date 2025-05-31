package user

import (
	"github.com/nogavadu/articles-service/internal/service"
)

type Implementation struct {
	userServ service.UserService
}

func New(userServ service.UserService) *Implementation {
	return &Implementation{
		userServ: userServ,
	}
}
