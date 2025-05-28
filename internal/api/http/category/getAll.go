package category

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
	"strconv"
)

type getAllResponse struct {
	Data []model.Category `json:"data"`
}

func (i *Implementation) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := categoryGetAllQueryParams(r)
		if err != nil {
			response.Err(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		categories, err := i.categoryServ.GetAll(r.Context(), params)
		if err != nil {
			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &getAllResponse{
			Data: categories,
		})
	}
}

func categoryGetAllQueryParams(r *http.Request) (*model.CategoryGetAllParams, error) {
	params := &model.CategoryGetAllParams{}

	cropIdStr := r.URL.Query().Get("crop_id")
	if cropIdStr != "" {
		id, err := strconv.Atoi(cropIdStr)
		if err != nil {
			return nil, errors.New("invalid crop_id query param")
		}
		params.CropId = &id
	}

	status := r.URL.Query().Get("status")
	if status != "" {
		params.Status = &status
	}

	return params, nil
}
