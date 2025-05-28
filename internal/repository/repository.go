package repository

import (
	"context"
	articleRepoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
	categoryRepoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
	cropRepoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
	statusRepoModel "github.com/nogavadu/articles-service/internal/repository/status/model"
)

type CropRepository interface {
	Create(ctx context.Context, info *cropRepoModel.CropInfo) (int, error)
	GetAll(ctx context.Context, statusId int) ([]cropRepoModel.Crop, error)
	GetById(ctx context.Context, id int) (*cropRepoModel.Crop, error)
	Update(ctx context.Context, id int, input *cropRepoModel.UpdateInput) error
	Delete(ctx context.Context, id int) error
}

type CategoryRepository interface {
	Create(ctx context.Context, info *categoryRepoModel.CategoryInfo) (int, error)
	GetAll(ctx context.Context, params *categoryRepoModel.CategoryGetAllParams) ([]categoryRepoModel.Category, error)
	GetById(ctx context.Context, id int) (*categoryRepoModel.Category, error)
	Update(ctx context.Context, id int, input *categoryRepoModel.UpdateInput) error
	Delete(ctx context.Context, id int) error
}

type CropCategoriesRepository interface {
	Create(ctx context.Context, cropId int, categoryId int) error
	Delete(ctx context.Context, cropId int, categoryId int) error
}

type ArticleRepository interface {
	Create(ctx context.Context, articleBody *articleRepoModel.ArticleBody) (int, error)
	GetAll(ctx context.Context, params *articleRepoModel.ArticleGetAllParams) ([]articleRepoModel.Article, error)
	GetById(ctx context.Context, id int) (*articleRepoModel.Article, error)
	Update(ctx context.Context, id int, input *articleRepoModel.UpdateInput) error
	Delete(ctx context.Context, id int) error
}

type ArticleRelationsRepository interface {
	Create(ctx context.Context, cropId int, categoryId int, articleId int) error
}

type ArticleImagesRepository interface {
	CreateBulk(ctx context.Context, articleId int, images []string) error
	GetAll(ctx context.Context, articleId int) ([]string, error)
	DeleteBulk(ctx context.Context, articleId int) error
}

type StatusRepository interface {
	Create(ctx context.Context, status string) (int, error)
	GetAll(ctx context.Context) ([]statusRepoModel.Status, error)
	GetByStatus(ctx context.Context, status string) (*statusRepoModel.Status, error)
}
