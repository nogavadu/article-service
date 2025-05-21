package model

import "time"

type ArticleGetAllParams struct {
	CropId     *int
	CategoryId *int
	Limit      *int
	Offset     *int
}

type Article struct {
	Id int `json:"id"`
	ArticleBody
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ArticleBody struct {
	Title  string   `json:"title" validate:"required"`
	Text   string   `json:"text,omitempty"`
	Images []string `json:"images,omitempty"`
}

type ArticleUpdateInput struct {
	Title  string   `json:"title,omitempty"`
	Text   string   `json:"text,omitempty"`
	Images []string `json:"images,omitempty"`
}
