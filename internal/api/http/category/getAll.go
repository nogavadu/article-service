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
	options := &model.CategoryGetAllParams{}

	cropIdStr := r.URL.Query().Get("crop_id")
	if cropIdStr != "" {
		id, err := strconv.Atoi(cropIdStr)
		if err != nil {
			return nil, errors.New("invalid crop_id query param")
		}
		options.CropId = &id
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, errors.New("invalid limit query param")
		}
		options.Limit = &l
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, errors.New("invalid offset query param")
		}
		options.Offset = &o
	}

	return options, nil
}
