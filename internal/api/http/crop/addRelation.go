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

type addRelationResponse struct {
	Status string `json:"status"`
}

func (i *Implementation) AddRelationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cropIdStr := chi.URLParam(r, "cropId")
		if cropIdStr == "" {
			response.Err(w, r, "crop id is required", http.StatusBadRequest)
			return
		}
		cropId, err := strconv.Atoi(cropIdStr)
		if err != nil {
			response.Err(w, r, "invalid crop id", http.StatusBadRequest)
			return
		}

		categoryIdStr := chi.URLParam(r, "categoryId")
		if categoryIdStr == "" {
			response.Err(w, r, "category id is required", http.StatusBadRequest)
			return
		}
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			response.Err(w, r, "invalid category id", http.StatusBadRequest)
		}

		if err = i.cropServ.AddRelation(r.Context(), cropId, categoryId); err != nil {
			if errors.Is(err, crop.ErrAccessDenied) {
				render.JSON(w, r, &updateResponse{
					Status: "AccessDenied",
				})
				return
			}

			response.Err(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		render.JSON(w, r, &addRelationResponse{
			Status: "ok",
		})
	}
}
