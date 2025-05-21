package article

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/request"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	"net/http"
	"strconv"
)

type UpdateRequest struct {
	model.ArticleUpdateInput
}

type UpdateResponse struct {
	Status string `json:"status"`
}

func (i *Implementation) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "articleId")
		if idStr == "" {
			response.Err(w, r, "article id is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Err(w, r, "invalid article id", http.StatusBadRequest)
			return
		}

		var reqData UpdateRequest
		if err = json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			response.Err(w, r, "invalid request body", http.StatusBadRequest)
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

		if err = i.articleServ.Update(r.Context(), id, &reqData.ArticleUpdateInput); err != nil {
			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, &UpdateResponse{
			Status: "ok",
		})
	}
}
