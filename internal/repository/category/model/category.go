package model

type Category struct {
	ID   int `db:"id"`
	Info CategoryInfo
}

type CategoryInfo struct {
	Name string `db:"name"`
}
