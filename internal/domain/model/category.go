package model

type Category struct {
	ID int `json:"id"`
	CategoryInfo
}

type CategoryInfo struct {
	Name string `json:"name" validate:"required"`
}
