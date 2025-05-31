package crop

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	cropServ "github.com/nogavadu/articles-service/internal/service/crop"
	"net/http"
)

type createRequest struct {
	UserId int            `json:"user_id" validate:"required"`
	Crop   model.CropInfo `json:"crop" validate:"required"`
}

type createResponse struct {
	Id int `json:"id"`
}

func (i *Implementation) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqData createRequest
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			response.Err(w, r, fmt.Sprintf("invalid request body: %s", err), http.StatusBadRequest)
			return
		}
		if err := validator.New().Struct(&reqData); err != nil {
			response.Err(w, r, fmt.Sprintf("invalid arguments: %s", err), http.StatusBadRequest)
			return
		}
		if reqData.Crop.Img != nil {
			if err := validator.New().Var(reqData.Crop.Img, "url"); err != nil {
				response.Err(w, r, "invalid image url", http.StatusBadRequest)
				return
			}
		}

		id, err := i.cropServ.Create(r.Context(), reqData.UserId, &reqData.Crop)
		if err != nil {
			if errors.Is(err, cropServ.ErrAlreadyExists) {
				response.Err(w, r, err.Error(), http.StatusBadRequest)
				return
			}
			if errors.Is(err, cropServ.ErrAccessDenied) {
				render.JSON(w, r, &updateResponse{
					Status: "AccessDenied",
				})
				return
			}

			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, &createResponse{
			Id: id,
		})
	}
}
