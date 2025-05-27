package auth

import (
	"github.com/go-chi/render"
	"net/http"
)

type getRefreshTokenResponse struct {
	Token string `json:"token"`
}

func (i *Implementation) GetRefreshTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := i.authServ.GetRefreshToken(r.Context())
		if err != nil {
			render.JSON(w, r, &getRefreshTokenResponse{
				Token: "",
			})
		}

		render.JSON(w, r, &getRefreshTokenResponse{
			Token: token,
		})
	}
}
