package model

import "time"

type ArticleGetAllParams struct {
	CropId     *int
	CategoryId *int
	Status     *string
}

type Article struct {
	Id int `json:"id"`
	ArticleBody
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticleBody struct {
	Title     string   `json:"title" validate:"required"`
	LatinName *string  `json:"latin_name,omitempty"`
	Text      *string  `json:"text,omitempty"`
	Images    []string `json:"images,omitempty"`
	Status    string   `json:"status"`
}

type ArticleUpdateInput struct {
	Title     *string  `json:"title,omitempty"`
	LatinName *string  `json:"latin_name,omitempty"`
	Text      *string  `json:"text,omitempty"`
	Images    []string `json:"images,omitempty"`
	Status    *string  `json:"status"`
}
