package model

type Crop struct {
	ID   int       `json:"id"`
	Info *CropInfo `json:"info"`
}

type CropInfo struct {
	Name string `json:"name" validate:"required"`
}
