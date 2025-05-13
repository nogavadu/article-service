package model

type Article struct {
	ID   uint64 `db:"id"`
	Body ArticleBody
}

type ArticleBody struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
