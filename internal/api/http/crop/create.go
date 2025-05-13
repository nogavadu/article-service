package crop

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	cropServ "github.com/nogavadu/articles-service/internal/service/crop"
	"net/http"
)

type createRequest struct {
	*model.CropInfo `json:"crop_info" validate:"required"`
}

type createResponse struct {
	Id int `json:"id"`
}

func (i *Implementation) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var reqData createRequest
		if json.NewDecoder(r.Body).Decode(&reqData) != nil {
			response.Err(w, r, "invalid request body format", http.StatusBadRequest)
			return
		}
		if err := validator.New().Struct(&reqData); err != nil {
			response.Err(w, r, "invalid arguments", http.StatusBadRequest)
			return
		}

		id, err := i.cropServ.Create(r.Context(), reqData.CropInfo)
		if err != nil {
			if errors.Is(err, cropServ.ErrAlreadyExists) {
				response.Err(w, r, err.Error(), http.StatusBadRequest)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, &createResponse{
			Id: id,
		})
	}
}
