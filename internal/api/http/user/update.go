package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/request"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
	"strconv"
)

type updateRequest struct {
	model.UserUpdateInput
}

type updateResponse struct {
	Status string `json:"status"`
}

func (i *Implementation) UpdateHandler() http.HandlerFunc {
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

		var reqBody updateRequest
		if err = json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			response.Err(w, r, fmt.Sprintf("invalid request body: %s", err), http.StatusBadRequest)
			return
		}
		isEmpty, err := request.IsStructEmpty(reqBody)
		if err != nil {
			response.Err(w, r, "invalid request body type", http.StatusBadRequest)
			return
		}
		if isEmpty {
			response.Err(w, r, "empty request body", http.StatusBadRequest)
			return
		}

		if err = i.userServ.Update(r.Context(), userId, &reqBody.UserUpdateInput); err != nil {
			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &updateResponse{
			Status: "ok",
		})
	}
}
