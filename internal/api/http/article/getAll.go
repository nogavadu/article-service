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

type getListResponse struct {
	Data []*model.Article `json:"data"`
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

		render.JSON(w, r, getListResponse{
			Data: articles,
		})
	}
}

func articleGetAllQueryParams(r *http.Request) (*model.ArticleGetAllParams, error) {
	options := &model.ArticleGetAllParams{}

	cropIdStr := r.URL.Query().Get("crop_id")
	if cropIdStr != "" {
		id, err := strconv.Atoi(cropIdStr)
		if err != nil {
			return nil, errors.New("invalid crop_id query param")
		}
		options.CropId = &id
	}

	categoryIdStr := r.URL.Query().Get("category_id")
	if categoryIdStr != "" {
		id, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			return nil, errors.New("invalid category_id query param")
		}
		options.CategoryId = &id
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
