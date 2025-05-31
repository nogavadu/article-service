package converter

import (
	"github.com/nogavadu/articles-service/internal/domain/model"
	repoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
)

func ToCrop(crop *repoModel.Crop, status string, author *model.User) *model.Crop {
	return &model.Crop{
		ID:        crop.ID,
		CropInfo:  *ToCropInfo(&crop.CropInfo, status, author),
		CreatedAt: crop.CreatedAt,
		UpdatedAt: crop.UpdatedAt,
	}
}

func ToCropInfo(cropInfo *repoModel.CropInfo, status string, author *model.User) *model.CropInfo {
	return &model.CropInfo{
		Name:        cropInfo.Name,
		Description: cropInfo.Description,
		Img:         cropInfo.Img,
		Status:      status,
		Author:      author,
	}
}

func ToRepoCropInfo(info *model.CropInfo, statusId int, authorId int) *repoModel.CropInfo {
	return &repoModel.CropInfo{
		Name:        info.Name,
		Description: info.Description,
		Img:         info.Img,
		Status:      statusId,
		Author:      &authorId,
	}
}

func ToRepoCropUpdateInput(input *model.UpdateCropInput, statusId *int) *repoModel.UpdateInput {
	return &repoModel.UpdateInput{
		Name:        input.Name,
		Description: input.Description,
		Img:         input.Img,
		Status:      statusId,
	}
}
