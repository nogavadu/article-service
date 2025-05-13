package article

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/response"
	articleService "github.com/nogavadu/articles-service/internal/service/article"
	"net/http"
	"strconv"
)

type getListResponse struct {
	Data []*model.Article `json:"data"`
}

func (i *Implementation) GetListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cropIdStr := chi.URLParam(r, "crop_id")
		if cropIdStr == "" {
			response.Err(w, r, "crop id is required", http.StatusBadRequest)
			return
		}
		cropId, err := strconv.Atoi(cropIdStr)
		if err != nil {
			response.Err(w, r, "invalid crop id", http.StatusBadRequest)
			return
		}

		categoryIdStr := chi.URLParam(r, "category_id")
		if categoryIdStr == "" {
			response.Err(w, r, "category id is required", http.StatusBadRequest)
			return
		}
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			response.Err(w, r, "invalid category id", http.StatusBadRequest)
			return
		}

		articles, err := i.articleServ.GetList(r.Context(), cropId, categoryId)
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
