package model

import "time"

type ArticleGetAllParams struct {
	CropId     *int
	CategoryId *int
	Limit      *int
	Offset     *int
}

type Article struct {
	Id int `db:"id"`
	ArticleBody
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ArticleBody struct {
	Title     string  `db:"title"`
	LatinName *string `db:"latin_name"`
	Text      *string `db:"text"`
}

type UpdateInput struct {
	Title     *string `db:"title"`
	LatinName *string `db:"latin_name"`
	Text      *string `db:"text"`
}
