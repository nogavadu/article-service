package category

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"io"
	"net/http"
)

type createRequest struct {
	Data *model.CategoryInfo `json:"data" validate:"required"`
}

type createResponse struct {
	Id int `json:"id"`
}

func (i *Implementation) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var reqData createRequest
		if err = json.Unmarshal(body, &reqData); err != nil {
			http.Error(w, "invalid request format", http.StatusBadRequest)
			return
		}

		if err = validator.New().Struct(&reqData); err != nil {
			http.Error(w, "invalid arguments", http.StatusBadRequest)
			return
		}

		id, err := i.categoryServ.Create(r.Context(), reqData.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		render.JSON(w, r, createResponse{
			Id: id,
		})
	}
}
