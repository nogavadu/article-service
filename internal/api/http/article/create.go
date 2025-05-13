package article

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"io"
	"net/http"
)

type createRequest struct {
	CropID      int               `json:"crop_id" validate:"required"`
	CategoryID  int               `json:"category_id" validate:"required"`
	ArticleBody model.ArticleBody `json:"article_body" validate:"required"`
}

type createResponse struct {
	Id int `json:"article_id"`
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

		id, err := i.articleServ.Create(r.Context(), reqData.CropID, reqData.CategoryID, &reqData.ArticleBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, &createResponse{
			Id: id,
		})
	}
}
