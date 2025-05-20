package model

type CategoryGetAllParams struct {
	CropId *int
	Limit  *int
	Offset *int
}

type Category struct {
	ID int `json:"id"`
	CategoryInfo
}

type CategoryInfo struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

type UpdateCategoryInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}
