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

func (r *articleRepository) Create(ctx context.Context, cropId int, categoryId int, article *articleRepoModel.ArticleBody) (int, error) {
	query := `
		INSERT INTO articles(crop_id, category_id, title, text)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var articleId int
	if err := r.db.QueryRow(ctx, query,
		cropId, categoryId, article.Title, article.Text,
	).Scan(&articleId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return 0, fmt.Errorf("%w: %w", ErrAlreadyExists, err)
			}
			if pgErr.Code == postgresErrors.InvalidForeignKeyErrCode {
				return 0, fmt.Errorf("%w: %w", ErrInvalidArguments, err)
			}
		}

		return 0, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return articleId, nil
}

func (r *articleRepository) GetById(ctx context.Context, id int) (*articleRepoModel.Article, error) {
	query := `
		SELECT id, title, text
		FROM articles
		WHERE id = $1;
	`

	var article articleRepoModel.Article
	if err := r.db.QueryRow(ctx, query, id).Scan(&article.ID, &article.Title, &article.Text); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return &article, nil
}

func (r *articleRepository) GetList(ctx context.Context, cropId int, categoryId int) ([]*articleRepoModel.Article, error) {
	query := `
		SELECT id, title, text
		FROM articles
		WHERE crop_id = $1 AND category_id = $2;
	`

	rows, err := r.db.Query(ctx, query, cropId, categoryId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}
	defer rows.Close()

	var articles []*articleRepoModel.Article
	for rows.Next() {
		var article articleRepoModel.Article

		if err = rows.Scan(&article.ID, &article.Title, &article.Text); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
		}

		articles = append(articles, &article)
	}

	return articles, nil
}
