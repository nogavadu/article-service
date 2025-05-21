package article

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nogavadu/articles-service/internal/client/db"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	articleRepoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
	"time"
)

var (
	ErrAlreadyExists       = errors.New("article already exists")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type articleRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repository.ArticleRepository {
	return &articleRepository{
		dbc: dbc,
	}
}

func (r *articleRepository) Create(
	ctx context.Context,
	articleBody *articleRepoModel.ArticleBody,
) (int, error) {
	queryRaw, args, err := sq.
		Insert("articles").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "text", "created_at", "updated_at").
		Values(articleBody.Title, articleBody.Text, time.Now(), time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to build query: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "articleRepository.Create",
		QueryRaw: queryRaw,
	}

	var articleId int
	if err = r.dbc.DB().ScanOneContext(ctx, &articleId, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return 0, fmt.Errorf("%w: %w", ErrAlreadyExists, err)
			}
		}

		return 0, fmt.Errorf("%w: failed to create article: %w", ErrInternalServerError, err)
	}

	return articleId, nil
}

func (r *articleRepository) GetAll(
	ctx context.Context,
	params *articleRepoModel.ArticleGetAllParams,
) ([]articleRepoModel.Article, error) {
	builder := sq.
		Select("a.id, a.title, a.text").
		PlaceholderFormat(sq.Dollar).
		From("articles AS a")

	if params.CropId != nil && params.CategoryId != nil {
		builder = builder.Join(
			"articles_relations AS ar ON a.id = ar.article_id AND ar.crop_id = ? AND ar.category_id = ?",
			*params.CropId, *params.CategoryId,
		)
	} else if params.CropId != nil {
		builder = builder.Join(
			"articles_relations AS ar ON a.id = ar.article_id AND ar.crop_id = ?",
			*params.CropId,
		)
	} else if params.CategoryId != nil {
		builder = builder.Join(
			"articles_relations AS ar ON a.id = ar.article_id AND ar.category_id = ?",
			*params.CategoryId,
		)
	}

	if params.Limit != nil {
		builder = builder.Limit(uint64(*params.Limit))
	}

	if params.Offset != nil {
		builder = builder.Offset(uint64(*params.Offset))
	}

	queryRaw, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "articleRepository.GetAll",
		QueryRaw: queryRaw,
	}

	var articles []articleRepoModel.Article
	if err = r.dbc.DB().ScanAllContext(ctx, &articles, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get articles: %s: %w", ErrInternalServerError, err)
	}

	return articles, nil
}

func (r *articleRepository) GetById(ctx context.Context, id int) (*articleRepoModel.Article, error) {
	queryRaw, args, err := sq.
		Select("id, title, text").
		PlaceholderFormat(sq.Dollar).
		From("articles").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "articleRepository.GetById",
		QueryRaw: queryRaw,
	}

	var article articleRepoModel.Article
	if err = r.dbc.DB().ScanOneContext(ctx, &article, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get article by id: %s: %w", ErrInternalServerError, err)
	}

	return &article, nil
}
