package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
)

func ToCrop(crop *repoModel.Crop) *model.Crop {
	return &model.Crop{
		ID:        crop.ID,
		CropInfo:  *ToCropInfo(crop.Info),
		CreatedAt: crop.CreatedAt,
		UpdatedAt: crop.UpdatedAt,
	}
}

func ToCropInfo(cropInfo *repoModel.CropInfo) *model.CropInfo {
	return &model.CropInfo{
		Name:        cropInfo.Name,
		Description: cropInfo.Description,
		Img:         cropInfo.Img,
	}
}

func ToRepoCropInfo(info *model.CropInfo, status int) *repoModel.CropInfo {
	return &repoModel.CropInfo{
		Name:        info.Name,
		Description: info.Description,
		Img:         info.Img,
		Status:      status,
	}
}

func ToRepoCropUpdateInput(input *model.UpdateCropInput) *repoModel.UpdateInput {
	return &repoModel.UpdateInput{
		Name:        input.Name,
		Description: input.Description,
		Img:         input.Img,
	}
}
