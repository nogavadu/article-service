package model

type Category struct {
	ID int `db:"id"`
	CategoryInfo
}

type CategoryInfo struct {
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Icon        *string `db:"icon"`
}
