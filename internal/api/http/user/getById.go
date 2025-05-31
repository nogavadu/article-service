package user

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
	"strconv"
)

type getByIdResponse struct {
	*model.User
}

func (i *Implementation) GetByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := chi.URLParam(r, "userId")
		if userIdStr == "" {
			response.Err(w, r, "user id is required", http.StatusBadRequest)
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			response.Err(w, r, "invalid user id", http.StatusBadRequest)
			return
		}

		user, err := i.userServ.GetById(r.Context(), userId)
		if err != nil {
			response.Err(w, r, err.Error(), http.StatusNotFound)
		}

		render.JSON(w, r, &getByIdResponse{
			User: user,
		})
	}
}
