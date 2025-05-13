package article

import (
	"context"
	"errors"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	repo "github.com/nogavadu/articles-service/internal/repository/article"
	"github.com/nogavadu/articles-service/internal/service"
	"log/slog"
)

var (
	ErrNotFound            = errors.New("article not found")
	ErrAlreadyExists       = errors.New("article already exists")
	ErrInvalidArguments    = errors.New("invalid article arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type articleService struct {
	log *slog.Logger

	articleRepo repository.ArticleRepository
}

func New(log *slog.Logger, articleRepository repository.ArticleRepository) service.ArticleService {
	return &articleService{
		log:         log,
		articleRepo: articleRepository,
	}
}

func (s *articleService) Create(ctx context.Context, cropID int, categoryID int, articleBody *model.ArticleBody) (int, error) {
	const op = "articleService.Create"
	log := s.log.With(slog.String("op", op))

	articleID, err := s.articleRepo.Create(ctx, cropID, categoryID, converter.ToRepoArticleBody(articleBody))
	if err != nil {
		if errors.Is(err, repo.ErrInvalidArguments) {
			return 0, ErrInvalidArguments
		}
		if errors.Is(err, repo.ErrAlreadyExists) {
			return 0, ErrAlreadyExists
		}

		log.Error("failed to create article", slog.String("error", err.Error()))
		return 0, ErrInternalServerError
	}

	return articleID, nil
}

func (s *articleService) GetByID(ctx context.Context, id int) (*model.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (s *articleService) GetList(ctx context.Context, cropID int, categoryID int) ([]*model.Article, error) {
	//TODO implement me
	panic("implement me")
}
