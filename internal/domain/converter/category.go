package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
)

func ToCategory(category *repoModel.Category, status string, author *model.User) *model.Category {
	return &model.Category{
		ID:           category.ID,
		CategoryInfo: *ToCategoryInfo(category, status, author),
	}
}

func ToCategoryInfo(category *repoModel.Category, status string, author *model.User) *model.CategoryInfo {
	return &model.CategoryInfo{
		Name:        category.Name,
		Description: category.Description,
		Icon:        category.Icon,
		Status:      status,
		Author:      author,
	}
}

func ToRepoCategoryInfo(categoryInfo *model.CategoryInfo, status int, author int) *repoModel.CategoryInfo {
	return &repoModel.CategoryInfo{
		Name:        categoryInfo.Name,
		Description: categoryInfo.Description,
		Status:      status,
		Author:      &author,
	}
}

func ToRepoCategoryGetAllParams(params *model.CategoryGetAllParams, status int) *repoModel.CategoryGetAllParams {
	return &repoModel.CategoryGetAllParams{
		CropId: params.CropId,
		Status: status,
	}
}

func ToRepoCategoryUpdateInput(input *model.UpdateCategoryInput, statusId *int) *repoModel.UpdateInput {
	return &repoModel.UpdateInput{
		Name:        input.Name,
		Description: input.Description,
		Icon:        input.Icon,
		Status:      statusId,
	}
}
