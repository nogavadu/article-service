package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
)

func ToRepoCropInfo(info *model.CropInfo) *repoModel.CropInfo {
	return &repoModel.CropInfo{
		Name: info.Name,
	}
}

func ToCrop(crop *repoModel.Crop) *model.Crop {
	return &model.Crop{
		ID:   crop.ID,
		Info: ToCropInfo(crop.Info),
	}
}

func ToCropInfo(cropInfo *repoModel.CropInfo) *model.CropInfo {
	return &model.CropInfo{
		Name: cropInfo.Name,
	}
}
