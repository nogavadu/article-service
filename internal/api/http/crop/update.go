package crop

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/request"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"github.com/nogavadu/articles-service/internal/service/crop"
	"net/http"
	"strconv"
)

type updateRequest struct {
	model.UpdateCropInput
}

type updateResponse struct {
	Status string `json:"status"`
}

func (i *Implementation) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "cropId")
		if idStr == "" {
			response.Err(w, r, "crop id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Err(w, r, "invalid crop id", http.StatusBadRequest)
			return
		}

		var reqData updateRequest
		if err = json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			response.Err(w, r, fmt.Sprintf("invalid request body: %s", err), http.StatusBadRequest)
			return
		}

		isEmpty, err := request.IsStructEmpty(reqData)
		if err != nil {
			response.Err(w, r, "invalid request body type", http.StatusBadRequest)
			return
		}
		if isEmpty {
			response.Err(w, r, "empty request body", http.StatusBadRequest)
			return
		}

		if err = i.cropServ.Update(r.Context(), id, &reqData.UpdateCropInput); err != nil {
			if errors.Is(err, crop.ErrAccessDenied) {
				render.JSON(w, r, &updateResponse{
					Status: "AccessDenied",
				})
				return
			}

			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &updateResponse{
			Status: "ok",
		})
	}
}
