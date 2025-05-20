package category

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

type UpdateRequest struct {
	model.UpdateCategoryInput
}

type updateResponse struct {
	Status string `json:"status"`
}

func (i *Implementation) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "categoryId")
		if idStr == "" {
			response.Err(w, r, "category id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Err(w, r, "category id is invalid", http.StatusBadRequest)
			return
		}

		var reqData UpdateRequest
		if err = json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			response.Err(w, r, fmt.Sprintf("invalid request body: %s", err), http.StatusBadRequest)
			return
		}

		isEmpty, err := request.IsStructEmpty(reqData.UpdateCategoryInput)
		if err != nil {
			response.Err(w, r, "invalid request body type", http.StatusBadRequest)
			return
		}
		if isEmpty {
			response.Err(w, r, "empty request body", http.StatusBadRequest)
			return
		}

		if err = i.categoryServ.Update(r.Context(), id, &reqData.UpdateCategoryInput); err != nil {
			response.Err(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		render.JSON(w, r, &updateResponse{
			Status: "ok",
		})
	}
}
