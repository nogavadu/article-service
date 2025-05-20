package crop

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
	"reflect"
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
		if isReqDataEmpty(&reqData.UpdateCropInput) {
			response.Err(w, r, "empty request body", http.StatusBadRequest)
			return
		}

		if err := i.cropServ.Update(r.Context(), id, &reqData.UpdateCropInput); err != nil {
			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &updateResponse{
			Status: "ok",
		})
	}
}

func isReqDataEmpty(reqData *model.UpdateCropInput) bool {
	val := reflect.ValueOf(reqData).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsNil() {
			return false
		}
	}
	return true
}
