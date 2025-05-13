package category

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	categoryServ "github.com/nogavadu/articles-service/internal/service/category"
	"net/http"
)

type createRequest struct {
	*model.CategoryInfo `json:"data" validate:"required"`
}

type createResponse struct {
	Id int `json:"id"`
}

func (i *Implementation) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqData createRequest
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			response.Err(w, r, "invalid request body format", http.StatusBadRequest)
			return
		}
		if err := validator.New().Struct(&reqData); err != nil {
			response.Err(w, r, "invalid arguments", http.StatusBadRequest)
			return
		}

		id, err := i.categoryServ.Create(r.Context(), reqData.CategoryInfo)
		if err != nil {
			if errors.Is(err, categoryServ.ErrAlreadyExists) {
				response.Err(w, r, err.Error(), http.StatusBadRequest)
			}

			response.Err(w, r, "internal server error", http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createResponse{
			Id: id,
		})
	}
}
