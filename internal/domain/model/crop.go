package model

import "time"

type CropGetAllParams struct {
	Status *string
}

type Crop struct {
	ID int `json:"id"`
	CropInfo
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CropInfo struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Img         *string `json:"img,omitempty"`
	Status      string  `json:"status" validate:"required"`
	Author      *User   `json:"author,omitempty"`
}

type UpdateCropInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Img         *string `json:"img,omitempty"`
	Status      *string `json:"status,omitempty"`
}
