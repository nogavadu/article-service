package service

import (
	"context"
	"github.com/nogavadu/articles-service/internal/domain/model"
)

type CropService interface {
	Create(ctx context.Context, cropInfo *model.CropInfo) (int, error)
	GetAll(ctx context.Context) ([]model.Crop, error)
	GetById(ctx context.Context, id int) (*model.Crop, error)
	Update(ctx context.Context, id int, input *model.UpdateCropInput) error
}

type CategoryService interface {
	Create(ctx context.Context, category *model.CategoryInfo) (int, error)
	GetAll(ctx context.Context, params *model.CategoryGetAllParams) ([]model.Category, error)
	GetById(ctx context.Context, id int) (*model.Category, error)
	Update(ctx context.Context, id int, input *model.UpdateCategoryInput) error
}

type ArticleService interface {
	Create(ctx context.Context, cropId int, categoryId int, articleBody *model.ArticleBody) (int, error)
	GetAll(ctx context.Context, params *model.ArticleGetAllParams) ([]*model.Article, error)
	GetById(ctx context.Context, id int) (*model.Article, error)
}
