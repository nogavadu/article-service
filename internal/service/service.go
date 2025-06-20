package service

import (
	"context"
	"github.com/nogavadu/articles-service/internal/domain/model"
)

type AuthService interface {
	Register(ctx context.Context, authData *model.UserRegisterData) (int, error)
	Login(ctx context.Context, authData *model.UserAuthData) (string, error)
	GetRefreshToken(ctx context.Context) (string, error)
}

type UserService interface {
	GetById(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, id int, input *model.UserUpdateInput) error
	Delete(ctx context.Context, id int) error
}

type CropService interface {
	Create(ctx context.Context, userId int, cropInfo *model.CropInfo) (int, error)
	GetAll(ctx context.Context, params *model.CropGetAllParams) ([]model.Crop, error)
	GetById(ctx context.Context, id int) (*model.Crop, error)
	Update(ctx context.Context, id int, input *model.UpdateCropInput) error
	Delete(ctx context.Context, id int) error

	AddRelation(ctx context.Context, cropId int, categoryId int) error
	RemoveRelation(ctx context.Context, cropId int, categoryId int) error
}

type CategoryService interface {
	Create(ctx context.Context, userId int, category *model.CategoryInfo, params *model.CategoryCreateParams) (int, error)
	GetAll(ctx context.Context, params *model.CategoryGetAllParams) ([]model.Category, error)
	GetById(ctx context.Context, id int) (*model.Category, error)
	Update(ctx context.Context, id int, input *model.UpdateCategoryInput) error
	Delete(ctx context.Context, id int) error
}

type ArticleService interface {
	Create(ctx context.Context, userId int, cropId int, categoryId int, articleBody *model.ArticleBody) (int, error)
	GetAll(ctx context.Context, params *model.ArticleGetAllParams) ([]model.Article, error)
	GetById(ctx context.Context, id int) (*model.Article, error)
	Update(ctx context.Context, id int, input *model.ArticleUpdateInput) error
	Delete(ctx context.Context, id int) error
}

type StatusService interface {
	GetByStatus(ctx context.Context, status string) (*model.Status, error)
}
