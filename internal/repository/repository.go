package repository

import (
	"context"
	"github.com/nogavadu/articles-service/internal/domain/model"
	articleRepoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
	categoryRepoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
	cropRepoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
)

type CropRepository interface {
	Create(ctx context.Context, info *cropRepoModel.CropInfo) (int, error)
	GetAll(ctx context.Context) ([]*cropRepoModel.Crop, error)
	GetById(ctx context.Context, id int) (*cropRepoModel.Crop, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, info *categoryRepoModel.CategoryInfo) (int, error)
	GetAll(ctx context.Context, params *model.CategoryGetAllParams) ([]*categoryRepoModel.Category, error)
	GetById(ctx context.Context, id int) (*categoryRepoModel.Category, error)
}

type ArticleRepository interface {
	Create(
		ctx context.Context,
		cropId int,
		categoryId int,
		articleBody *articleRepoModel.ArticleBody,
		images []string,
	) (int, error)
	GetAll(ctx context.Context, params *model.ArticleGetAllParams) ([]*articleRepoModel.Article, error)
	GetById(ctx context.Context, id int) (*articleRepoModel.Article, error)
}
