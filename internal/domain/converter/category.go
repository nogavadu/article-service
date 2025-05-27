package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
)

func ToCategory(category *repoModel.Category) *model.Category {
	return &model.Category{
		ID: category.ID,
		CategoryInfo: model.CategoryInfo{
			Name:        category.Name,
			Description: category.Description,
			Icon:        category.Icon,
		},
	}
}

func ToRepoCategoryInfo(categoryInfo *model.CategoryInfo) *repoModel.CategoryInfo {
	return &repoModel.CategoryInfo{
		Name:        categoryInfo.Name,
		Description: categoryInfo.Description,
	}
}

func ToRepoCategoryGetAllParams(params *model.CategoryGetAllParams) *repoModel.CategoryGetAllParams {
	return &repoModel.CategoryGetAllParams{
		CropId: params.CropId,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
}

func ToRepoCategoryUpdateInput(input *model.UpdateCategoryInput) *repoModel.UpdateInput {
	return &repoModel.UpdateInput{
		Name:        input.Name,
		Description: input.Description,
		Icon:        input.Icon,
	}
}
