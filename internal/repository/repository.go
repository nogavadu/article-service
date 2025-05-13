package repository

import (
	"context"
	articleRepoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
	categoryRepoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
	cropRepoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
)

type CropRepository interface {
	Create(ctx context.Context, info *cropRepoModel.CropInfo) (int, error)
	GetAll(ctx context.Context) ([]*cropRepoModel.Crop, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, info *categoryRepoModel.CategoryInfo) (int, error)
	GetList(ctx context.Context, cropId int) ([]*categoryRepoModel.Category, error)
	GetAll(ctx context.Context) ([]*categoryRepoModel.Category, error)
}

type ArticleRepository interface {
	Create(ctx context.Context, cropId int, categoryId int, article *articleRepoModel.ArticleBody) (int, error)
	GetById(ctx context.Context, id int) (*articleRepoModel.Article, error)
	GetList(ctx context.Context, cropId int, categoryId int) ([]*articleRepoModel.Article, error)
}
