package model

type Article struct {
	ID uint64 `db:"id"`
	ArticleBody
}

type ArticleBody struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
