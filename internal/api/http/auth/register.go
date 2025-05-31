package auth

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
)

type registerRequest struct {
	model.UserRegisterData
}

type registerResponse struct {
	UserId int `json:"user_id"`
}

func (i *Implementation) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqData registerRequest
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			response.Err(w, r, "invalid request", http.StatusBadRequest)
			return
		}
		if err := validator.New().Struct(reqData); err != nil {
			response.Err(w, r, "invalid email/ password", http.StatusBadRequest)
			return
		}

		userId, err := i.authServ.Register(r.Context(), &reqData.UserRegisterData)
		if err != nil {
			response.Err(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		render.JSON(w, r, &registerResponse{
			UserId: userId,
		})
	}
}
