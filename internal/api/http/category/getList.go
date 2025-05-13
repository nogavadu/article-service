package category

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"net/http"
	"strconv"
)

type getListResponse struct {
	Data []*model.Category `json:"data"`
}

func (i *Implementation) GetListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cropIdStr := chi.URLParam(r, "crop_id")
		if cropIdStr == "" {
			http.Error(w, "crop id is required", http.StatusBadRequest)
			return
		}
		cropId, err := strconv.Atoi(cropIdStr)
		if err != nil {
			http.Error(w, "crop id is invalid", http.StatusBadRequest)
		}

		categories, err := i.categoryServ.GetList(r.Context(), cropId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		render.JSON(w, r, &getListResponse{
			Data: categories,
		})
	}
}
