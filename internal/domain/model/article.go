package model

type ArticleGetAllParams struct {
	CropId     *int
	CategoryId *int
	Limit      *int
	Offset     *int
}

type Article struct {
	Id uint64 `json:"id"`
	ArticleBody
}

type ArticleBody struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"body,omitempty"`
}
