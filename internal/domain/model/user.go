package model

type User struct {
	Id int `json:"id"`
	UserAuthData
	Role string `json:"role"`
}

type UserAuthData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
