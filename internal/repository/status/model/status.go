package model

type Status struct {
	Id     int    `db:"id"`
	Status string `db:"status"`
}
