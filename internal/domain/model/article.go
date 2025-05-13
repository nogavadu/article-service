package model

type Article struct {
	ID   uint64       `json:"id"`
	Body *ArticleBody `json:"body"`
}

type ArticleBody struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"body,omitempty"`
}
