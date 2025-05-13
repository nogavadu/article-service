package article

import (
	"context"
	"errors"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	articleRepo "github.com/nogavadu/articles-service/internal/repository/article"
	"github.com/nogavadu/articles-service/internal/service"
	"log/slog"
)

var (
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

func (s *articleService) Create(ctx context.Context, cropId int, categoryId int, articleBody *model.ArticleBody) (int, error) {
	const op = "articleService.Create"
	log := s.log.With(slog.String("op", op))

	articleId, err := s.articleRepo.Create(ctx, cropId, categoryId, converter.ToRepoArticleBody(articleBody))
	if err != nil {
		log.Error("failed to create article", slog.String("error", err.Error()))

		if errors.Is(err, articleRepo.ErrInvalidArguments) {
			return 0, ErrInvalidArguments
		}
		if errors.Is(err, articleRepo.ErrAlreadyExists) {
			return 0, ErrAlreadyExists
		}

		return 0, ErrInternalServerError
	}

	return articleId, nil
}

func (s *articleService) GetAll(ctx context.Context, params *model.ArticleGetAllParams) ([]*model.Article, error) {
	const op = "articleService.GetAll"
	log := s.log.With(slog.String("op", op))

	repoArticles, err := s.articleRepo.GetAll(ctx, params)
	if err != nil {
		log.Error("failed to get articles", slog.String("error", err.Error()))

		if errors.Is(err, articleRepo.ErrInvalidArguments) {
			return nil, ErrInvalidArguments
		}

		return nil, ErrInternalServerError
	}

	articles := make([]*model.Article, 0, len(repoArticles))
	for _, a := range repoArticles {
		articles = append(articles, converter.ToArticle(a))
	}

	return articles, nil
}

func (s *articleService) GetById(ctx context.Context, id int) (*model.Article, error) {
	const op = "articleService.GetById"
	log := s.log.With(slog.String("op", op))

	article, err := s.articleRepo.GetById(ctx, id)
	if err != nil {
		log.Error("failed to get article", slog.String("error", err.Error()))

		return nil, ErrInternalServerError
	}

	return converter.ToArticle(article), nil
}
