package model

type Crop struct {
	ID   int `db:"id"`
	Info *CropInfo
}

type CropInfo struct {
	Name string `db:"name"`
}
