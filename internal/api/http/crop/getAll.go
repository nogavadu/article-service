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
		crops, err := i.cropServ.GetAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &getAllResponse{
			Data: crops,
		})
	}
}
