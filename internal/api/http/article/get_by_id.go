package article

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"net/http"
	"strconv"
)

type GetByIDResponse struct {
	model.Article
}

func (i *Implementation) GetByIDHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			http.Error(w, "empty article ID", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid article ID", http.StatusBadRequest)
			return
		}

		article, err := i.articleServ.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, GetByIDResponse{
			*article,
		})
	}
}
