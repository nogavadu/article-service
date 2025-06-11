package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
)

func ToArticle(article *repoModel.Article, images []string, status string, author *model.User) *model.Article {
	return &model.Article{
		Id:          article.Id,
		ArticleBody: *ToArticleBody(article, images, status, author),
	}
}

func ToArticleBody(article *repoModel.Article, images []string, status string, author *model.User) *model.ArticleBody {
	return &model.ArticleBody{
		Title:     article.Title,
		Text:      article.Text,
		LatinName: article.LatinName,
		Images:    images,
		Author:    author,
		Status:    status,
	}
}

func ToRepoArticleBody(body *model.ArticleBody, status int, author int) *repoModel.ArticleBody {
	return &repoModel.ArticleBody{
		Title:     body.Title,
		LatinName: body.LatinName,
		Text:      body.Text,
		Status:    status,
		Author:    &author,
	}
}

func ToRepoArticleGetAllParams(params *model.ArticleGetAllParams, status int) *repoModel.ArticleGetAllParams {
	return &repoModel.ArticleGetAllParams{
		CropId:     params.CropId,
		CategoryId: params.CategoryId,
		Status:     status,
	}
}

func ToRepoArticleUpdateInput(input *model.ArticleUpdateInput, statusId *int) *repoModel.UpdateInput {
	return &repoModel.UpdateInput{
		Title:     input.Title,
		LatinName: input.LatinName,
		Text:      input.Text,
		Status:    statusId,
	}
}
