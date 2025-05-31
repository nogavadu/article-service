package model

import "time"

type CategoryGetAllParams struct {
	CropId *int
	Status int
}

type Category struct {
	ID int `db:"id"`
	CategoryInfo
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CategoryInfo struct {
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Icon        *string `db:"icon"`
	Status      int     `db:"status"`
	Author      *int    `db:"author"`
}

type UpdateInput struct {
	Name        *string `db:"name"`
	Description *string `db:"description"`
	Icon        *string `db:"icon"`
	Status      *int    `db:"status"`
}
