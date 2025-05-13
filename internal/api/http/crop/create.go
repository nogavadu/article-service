package crop

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	cropServ "github.com/nogavadu/articles-service/internal/service/crop"
	"io"
	"net/http"
)

type createRequest struct {
	CropInfo *model.CropInfo `json:"crop_info" validate:"required"`
}

type createResponse struct {
	Id int `json:"crop_id"`
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

		id, err := i.cropServ.Create(r.Context(), reqData.CropInfo)
		if err != nil {
			if errors.Is(err, cropServ.ErrAlreadyExists) {
				http.Error(w, "crop already exists", http.StatusBadRequest)
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, &createResponse{
			Id: id,
		})
	}
}
