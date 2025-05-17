package crop

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
	"strconv"
)

type getByIdResponse struct {
	model.Crop
}

func (i *Implementation) GetByIdHandler() http.HandlerFunc {
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

		crop, err := i.cropServ.GetById(r.Context(), cropId)
		if err != nil {
			response.Err(w, r, err.Error(), http.StatusInternalServerError)
		}

		render.JSON(w, r, &getByIdResponse{
			Crop: *crop,
		})
	}
}
