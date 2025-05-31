package article

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/api/response"
	articleServ "github.com/nogavadu/articles-service/internal/service/article"
	"net/http"
)

type createRequest struct {
	UserId      int               `json:"user_id" validate:"required"`
	CropId      int               `json:"crop_id" validate:"required"`
	CategoryId  int               `json:"category_id" validate:"required"`
	ArticleBody model.ArticleBody `json:"article_body" validate:"required"`
}

type createResponse struct {
	Id int `json:"id"`
}

func (i *Implementation) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqData createRequest
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			response.Err(w, r, "invalid request body format", http.StatusBadRequest)
			return
		}
		if err := validator.New().Struct(&reqData); err != nil {
			response.Err(w, r, "invalid arguments", http.StatusBadRequest)
			return
		}

		id, err := i.articleServ.Create(r.Context(), reqData.UserId, reqData.CropId, reqData.CategoryId, &reqData.ArticleBody)
		if err != nil {
			if errors.Is(err, articleServ.ErrInvalidArguments) || errors.Is(err, articleServ.ErrAlreadyExists) {
				response.Err(w, r, err.Error(), http.StatusBadRequest)
				return
			}

			response.Err(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, &createResponse{
			Id: id,
		})
	}
}
