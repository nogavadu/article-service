package model

import "time"

type CategoryCreateParams struct {
	CropId *int
}

type CategoryGetAllParams struct {
	CropId *int
	Status *string
}

type Category struct {
	ID int `json:"id"`
	CategoryInfo
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CategoryInfo struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Status      string  `json:"status"`
}

type UpdateCategoryInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Status      *string `json:"status"`
}
