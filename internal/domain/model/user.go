package model

type User struct {
	Id int `json:"id"`
	UserInfo
}

type UserInfo struct {
	Name   *string `json:"name,omitempty"`
	Email  string  `json:"email"`
	Avatar *string `json:"avatar,omitempty"`
	Role   string  `json:"role"`
}

type UserRegisterData struct {
	Name string `json:"name,omitempty"`
	UserAuthData
}

type UserAuthData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateInput struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
	Role   *string `json:"role,omitempty"`
}
