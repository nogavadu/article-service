package category

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	categoryServ "github.com/nogavadu/articles-service/internal/service/category"
	"net/http"
	"strconv"
)

type getByIdResponse struct {
	*model.Category
}

func (i *Implementation) GetByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryIdStr := chi.URLParam(r, "categoryId")
		if categoryIdStr == "" {
			response.Err(w, r, "category id is required", http.StatusBadRequest)
			return
		}
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			response.Err(w, r, "invalid category id", http.StatusBadRequest)
			return
		}

		category, err := i.categoryServ.GetById(r.Context(), categoryId)
		if err != nil {
			if errors.Is(err, categoryServ.ErrNotFound) || errors.Is(err, categoryServ.ErrInvalidArguments) {
				response.Err(w, r, err.Error(), http.StatusBadRequest)
				return
			}

			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &getByIdResponse{
			Category: category,
		})
	}
}
