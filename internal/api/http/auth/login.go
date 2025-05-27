package auth

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
)

type loginRequest struct {
	model.UserAuthData
}

type loginResponse struct {
	Token string `json:"token"`
}

func (i *Implementation) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqData loginRequest
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			response.Err(w, r, "invalid request", http.StatusBadRequest)
			return
		}
		if err := validator.New().Struct(reqData); err != nil {
			response.Err(w, r, "invalid email/ password", http.StatusBadRequest)
			return
		}

		token, err := i.authServ.Login(r.Context(), &reqData.UserAuthData)
		if err != nil {
			response.Err(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		render.JSON(w, r, &loginResponse{
			Token: token,
		})
	}
}
