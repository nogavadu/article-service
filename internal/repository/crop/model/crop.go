package model

import "time"

type Crop struct {
	ID int `db:"id"`
	CropInfo
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CropInfo struct {
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Img         *string `db:"img"`
	Status      int     `db:"status"`
	Author      *int    `db:"author"`
}

type UpdateInput struct {
	Name        *string `db:"name"`
	Description *string `db:"description"`
	Img         *string `db:"img"`
	Status      *int    `db:"status"`
}
