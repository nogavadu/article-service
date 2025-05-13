package category

import (
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"net/http"
)

type getAllResponse struct {
	Data []*model.Category `json:"data"`
}

func (i *Implementation) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := i.categoryServ.GetAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		render.JSON(w, r, &getAllResponse{
			Data: categories,
		})
	}
}
