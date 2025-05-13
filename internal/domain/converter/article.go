package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
)

func ToRepoArticleBody(body *model.ArticleBody) *repoModel.ArticleBody {
	return &repoModel.ArticleBody{
		Title: body.Title,
		Text:  body.Text,
	}
}

func ToArticle(article *repoModel.Article) *model.Article {
	return &model.Article{
		Id: article.ID,
		ArticleBody: model.ArticleBody{
			Title: article.Title,
			Text:  article.Text,
		},
	}
}
