package article

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	articleRepoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
)

var (
	ErrAlreadyExists       = errors.New("article already exists")
	ErrNotFound            = errors.New("article not found")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type articleRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) Create(ctx context.Context, cropID int, categoryID int, article *articleRepoModel.ArticleBody) (int, error) {
	const op = "articleRepository.Create"

	query := `
		INSERT INTO articles(crop_id, category_id, title, text)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var articleId int
	if err := r.db.QueryRow(ctx, query,
		cropID, categoryID, article.Title, article.Text,
	).Scan(&articleId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return 0, fmt.Errorf("%s: %w: %w", op, ErrAlreadyExists, err)
			}
			if pgErr.Code == postgresErrors.InvalidForeignKeyErrCode {
				return 0, fmt.Errorf("%s: %w: %w", op, ErrInvalidArguments, err)
			}
		}

		return 0, fmt.Errorf("%s: %w: %w", op, ErrInternalServerError, err)
	}

	return articleId, nil
}

func (r *articleRepository) GetById(ctx context.Context, id int) (*articleRepoModel.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (r *articleRepository) GetList(ctx context.Context, cropID int, categoryID int) ([]*articleRepoModel.Article, error) {
	//TODO implement me
	panic("implement me")
}
