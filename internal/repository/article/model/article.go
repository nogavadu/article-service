package model

import "time"

type Article struct {
	Id int `db:"id"`
	ArticleBody
	Images    []string
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type ArticleBody struct {
	Title string `db:"title"`
	Text  string `db:"text"`
}
