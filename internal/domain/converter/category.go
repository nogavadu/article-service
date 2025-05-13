package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
)

func ToRepoCategoryInfo(categoryInfo *model.CategoryInfo) *repoModel.CategoryInfo {
	return &repoModel.CategoryInfo{
		Name: categoryInfo.Name,
	}
}

func ToCategory(category *repoModel.Category) *model.Category {
	return &model.Category{
		ID: category.ID,
		CategoryInfo: model.CategoryInfo{
			Name: category.Name,
		},
	}
}
