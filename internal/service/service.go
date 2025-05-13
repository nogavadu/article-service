package service

import (
	"context"
	"github.com/nogavadu/articles-service/internal/domain/model"
)

type CropService interface {
	Create(ctx context.Context, cropInfo *model.CropInfo) (int, error)
	GetAll(ctx context.Context) ([]*model.Crop, error)
}

type CategoryService interface {
	Create(ctx context.Context, category *model.CategoryInfo) (int, error)
	GetList(ctx context.Context, cropId int) ([]*model.Category, error)
	GetAll(ctx context.Context) ([]*model.Category, error)
}

type ArticleService interface {
	Create(ctx context.Context, cropId int, categoryId int, articleBody *model.ArticleBody) (int, error)
	GetById(ctx context.Context, id int) (*model.Article, error)
	GetList(ctx context.Context, cropId int, categoryId int) ([]*model.Article, error)
}
