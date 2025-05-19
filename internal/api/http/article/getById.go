package article

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	articleServ "github.com/nogavadu/articles-service/internal/service/article"
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
			response.Err(w, r, "article id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Err(w, r, "invalid article id", http.StatusBadRequest)
			return
		}

		article, err := i.articleServ.GetById(r.Context(), id)
		if err != nil {
			if errors.Is(err, articleServ.ErrInvalidArguments) {
				response.Err(w, r, err.Error(), http.StatusBadRequest)
				return
			}
			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, GetByIDResponse{
			Article: *article,
		})
	}
}
