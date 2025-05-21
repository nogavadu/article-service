package article_relations

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nogavadu/articles-service/internal/client/db"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
)

var (
	ErrAlreadyExists       = errors.New("article already exists")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type articleRelationsRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repository.ArticleRelationsRepository {
	return &articleRelationsRepository{
		dbc: dbc,
	}
}

func (r *articleRelationsRepository) Create(ctx context.Context, cropId int, categoryId int, articleId int) error {
	queryRaw, args, err := sq.
		Insert("articles_relations").
		PlaceholderFormat(sq.Dollar).
		Columns("crop_id", "category_id", "article_id").
		Values(cropId, categoryId, articleId).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %s: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "articleImagesRepository.CreateBulk",
		QueryRaw: queryRaw,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return fmt.Errorf("failed to create article relations: %w: %w", ErrAlreadyExists, err)
			}
			if pgErr.Code == postgresErrors.InvalidForeignKeyErrCode {
				return fmt.Errorf("failed to create article relations: %w: %w", ErrInvalidArguments, err)
			}
		}

		return fmt.Errorf("failed to create article relations: %s: %w", ErrInternalServerError, err)
	}

	return nil
}
