package crop

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"github.com/nogavadu/articles-service/internal/service/crop"
	"net/http"
	"strconv"
)

type DeleteResponse struct {
	Status string `json:"status"`
}

func (i *Implementation) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "cropId")
		if idStr == "" {
			response.Err(w, r, "article id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Err(w, r, "invalid article id", http.StatusBadRequest)
			return
		}

		err = i.cropServ.Delete(r.Context(), id)
		if err != nil {
			if errors.Is(err, crop.ErrAccessDenied) {
				render.JSON(w, r, &updateResponse{
					Status: "AccessDenied",
				})
				return
			}

			response.Err(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		render.JSON(w, r, DeleteResponse{Status: "ok"})
	}
}
