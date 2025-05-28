package crop

import (
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"net/http"
)

type getAllResponse struct {
	Data []model.Crop `json:"data"`
}

func (i *Implementation) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := cropGetAllParams(r)
		crops, err := i.cropServ.GetAll(r.Context(), params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &getAllResponse{
			Data: crops,
		})
	}
}

func cropGetAllParams(r *http.Request) *model.CropGetAllParams {
	params := &model.CropGetAllParams{}

	status := r.URL.Query().Get("status")
	if status != "" {
		params.Status = &status
	}

	return params
}
