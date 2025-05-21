package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
)

func ToArticle(article *repoModel.Article, images []string) *model.Article {
	return &model.Article{
		Id: article.Id,
		ArticleBody: model.ArticleBody{
			Title:  article.Title,
			Text:   article.Text,
			Images: images,
		},
	}
}

func ToRepoArticleBody(body *model.ArticleBody) *repoModel.ArticleBody {
	return &repoModel.ArticleBody{
		Title: body.Title,
		Text:  body.Text,
	}
}

func ToRepoArticleGetAllParams(params *model.ArticleGetAllParams) *repoModel.ArticleGetAllParams {
	return &repoModel.ArticleGetAllParams{
		CropId:     params.CropId,
		CategoryId: params.CategoryId,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}
}
