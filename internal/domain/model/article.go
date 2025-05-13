package model

type Article struct {
	Id uint64 `json:"id"`
	ArticleBody
}

type ArticleBody struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"body,omitempty"`
}
