package article

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"net/http"
	"strconv"
)

type getListResponse struct {
	Payload []*model.Article `json:"payload"`
}

func (i *Implementation) GetListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cropIdStr := chi.URLParam(r, "crop_id")
		if cropIdStr == "" {
			http.Error(w, "crop_id is required", http.StatusBadRequest)
			return
		}
		cropId, err := strconv.Atoi(cropIdStr)
		if err != nil {
			http.Error(w, "invalid crop id", http.StatusBadRequest)
		}

		categoryIdStr := chi.URLParam(r, "category_id")
		if categoryIdStr == "" {
			http.Error(w, "category_id is required", http.StatusBadRequest)
		}
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			http.Error(w, "invalid category id", http.StatusBadRequest)
		}

		articles, err := i.articleServ.GetList(r.Context(), cropId, categoryId)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, getListResponse{
			Payload: articles,
		})
	}
}
