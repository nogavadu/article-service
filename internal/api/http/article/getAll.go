package article

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	articleService "github.com/nogavadu/articles-service/internal/service/article"
	"net/http"
	"strconv"
)

type getAllResponse struct {
	Data []model.Article `json:"data"`
}

func (i *Implementation) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := articleGetAllQueryParams(r)
		if err != nil {
			response.Err(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		articles, err := i.articleServ.GetAll(r.Context(), params)
		if err != nil {
			if errors.Is(err, articleService.ErrInvalidArguments) {
				response.Err(w, r, err.Error(), http.StatusBadRequest)
				return
			}

			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &getAllResponse{
			Data: articles,
		})
	}
}

func articleGetAllQueryParams(r *http.Request) (*model.ArticleGetAllParams, error) {
	params := &model.ArticleGetAllParams{}

	cropIdStr := r.URL.Query().Get("crop_id")
	if cropIdStr != "" {
		id, err := strconv.Atoi(cropIdStr)
		if err != nil {
			return nil, errors.New("invalid crop_id query param")
		}
		params.CropId = &id
	}

	categoryIdStr := r.URL.Query().Get("category_id")
	if categoryIdStr != "" {
		id, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			return nil, errors.New("invalid category_id query param")
		}
		params.CategoryId = &id
	}

	status := r.URL.Query().Get("status")
	if status != "" {
		params.Status = &status
	}

	return params, nil
}
