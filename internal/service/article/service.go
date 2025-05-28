package article

import (
	"context"
	"errors"
	authService "github.com/nogavadu/articles-service/internal/clients/auth-service/grpc"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	articleRepo "github.com/nogavadu/articles-service/internal/repository/article"
	"github.com/nogavadu/articles-service/internal/service"
	"github.com/nogavadu/platform_common/pkg/db"
	"log/slog"
)

var (
	ErrAlreadyExists       = errors.New("article already exists")
	ErrInvalidArguments    = errors.New("invalid article arguments")
	ErrInternalServerError = errors.New("internal server error")
	ErrAccessDenied        = errors.New("access denied")
)

type articleService struct {
	log *slog.Logger

	articleRepo          repository.ArticleRepository
	articleImagesRepo    repository.ArticleImagesRepository
	articleRelationsRepo repository.ArticleRelationsRepository
	statusRepo           repository.StatusRepository

	txManager db.TxManager

	accessClient *authService.AccessServiceClient
	authClient   *authService.AuthServiceClient
}

func New(
	log *slog.Logger,
	articleRepository repository.ArticleRepository,
	articleImagesRepo repository.ArticleImagesRepository,
	articleRelationsRepo repository.ArticleRelationsRepository,
	statusRepo repository.StatusRepository,
	txManager db.TxManager,
	accessClient *authService.AccessServiceClient,
	authClient *authService.AuthServiceClient,
) service.ArticleService {
	return &articleService{
		log:                  log,
		articleRepo:          articleRepository,
		articleImagesRepo:    articleImagesRepo,
		articleRelationsRepo: articleRelationsRepo,
		statusRepo:           statusRepo,
		txManager:            txManager,
		accessClient:         accessClient,
		authClient:           authClient,
	}
}

func (s *articleService) Create(ctx context.Context, cropId int, categoryId int, articleBody *model.ArticleBody) (int, error) {
	const op = "articleService.Create"
	log := s.log.With(slog.String("op", op))

	var articleId int
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to create article", slog.String("error", errTx.Error()))
			}
		}()

		token, err := s.authClient.AccessToken(ctx)
		if err != nil {
			log.Error("failed to get access token", slog.String("error", err.Error()))
			return ErrAccessDenied
		}

		status, err := s.statusRepo.GetByStatus(ctx, articleBody.Status)
		if err != nil {
			log.Error("failed to get status", slog.String("error", err.Error()))
		}

		accessLevel := authService.ModeratorAccessLevel
		if status != nil && status.Id == 2 {
			accessLevel = authService.UserAccessLevel
		}

		err = s.accessClient.Check(ctx, token, accessLevel)
		if err != nil {
			log.Error("access check failed", slog.String("error", err.Error()))
			return ErrAccessDenied
		}

		articleId, errTx = s.articleRepo.Create(ctx, converter.ToRepoArticleBody(articleBody, status.Id))
		if errTx != nil {
			if errors.Is(errTx, articleRepo.ErrAlreadyExists) {
				return ErrAlreadyExists
			}

			return ErrInternalServerError
		}

		if len(articleBody.Images) > 0 {
			if errTx = s.articleImagesRepo.CreateBulk(ctx, articleId, articleBody.Images); errTx != nil {
				return ErrInternalServerError
			}
		}

		if errTx = s.articleRelationsRepo.Create(ctx, cropId, categoryId, articleId); errTx != nil {
			return ErrInternalServerError
		}

		return nil
	})

	return articleId, err
}

func (s *articleService) GetAll(ctx context.Context, params *model.ArticleGetAllParams) ([]model.Article, error) {
	const op = "articleService.GetAll"
	log := s.log.With(slog.String("op", op))

	var articles []model.Article
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to get articles", slog.String("error", errTx.Error()))
			}
		}()

		var statusId int
		if params.Status != nil {
			status, errTx := s.statusRepo.GetByStatus(ctx, *params.Status)
			if errTx != nil {
				log.Error("failed to get status", slog.String("error", errTx.Error()))
			}
			if status != nil {
				statusId = status.Id
			}
		} else {
			statusId = 2
		}

		repoArticles, errTx := s.articleRepo.GetAll(ctx, converter.ToRepoArticleGetAllParams(params, statusId))
		if errTx != nil {
			return ErrInternalServerError
		}

		articles = make([]model.Article, 0, len(repoArticles))
		for _, a := range repoArticles {
			imgs, errTx := s.articleImagesRepo.GetAll(ctx, a.Id)
			if errTx != nil {
				return ErrInternalServerError
			}
			articles = append(articles, *converter.ToArticle(&a, imgs))
		}

		return nil
	})

	return articles, err
}

func (s *articleService) GetById(ctx context.Context, id int) (*model.Article, error) {
	const op = "articleService.GetById"
	log := s.log.With(slog.String("op", op))

	var article *model.Article
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to get articles", slog.String("error", errTx.Error()))
			}
		}()

		repoArticle, errTx := s.articleRepo.GetById(ctx, id)
		if errTx != nil {
			return ErrInternalServerError
		}

		images, errTx := s.articleImagesRepo.GetAll(ctx, repoArticle.Id)
		if errTx != nil {
			return ErrInternalServerError
		}

		article = converter.ToArticle(repoArticle, images)

		return nil
	})

	return article, err
}

func (s *articleService) Update(ctx context.Context, id int, input *model.ArticleUpdateInput) error {
	const op = "articleService.Update"
	log := s.log.With(slog.String("op", op))

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to update article", slog.String("error", errTx.Error()))
			}
		}()

		errTx = s.articleRepo.Update(ctx, id, converter.ToRepoArticleUpdateInput(input))
		if errTx != nil {
			return ErrInternalServerError
		}

		errTx = s.articleImagesRepo.DeleteBulk(ctx, id)
		if errTx != nil {
			return ErrInternalServerError
		}

		if len(input.Images) > 0 {
			errTx = s.articleImagesRepo.CreateBulk(ctx, id, input.Images)
			if errTx != nil {
				return ErrInternalServerError
			}
		}

		return nil
	})

	return err
}

func (s *articleService) Delete(ctx context.Context, id int) error {
	const op = "articleService.Delete"
	log := s.log.With(slog.String("op", op))

	err := s.articleRepo.Delete(ctx, id)
	if err != nil {
		log.Error("failed to delete article", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}
