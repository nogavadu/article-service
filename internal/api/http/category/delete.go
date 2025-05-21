package category

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
	"strconv"
)

type DeleteResponse struct {
	Status string `json:"status"`
}

func (i *Implementation) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "categoryId")
		if idStr == "" {
			response.Err(w, r, "article id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Err(w, r, "invalid article id", http.StatusBadRequest)
			return
		}

		err = i.categoryServ.Delete(r.Context(), id)
		if err != nil {
			response.Err(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		render.JSON(w, r, DeleteResponse{Status: "ok"})
	}
}
