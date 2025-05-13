package category

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	categoryRepoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
	"strings"
)

var (
	ErrAlreadyExists       = errors.New("category already exists")
	ErrNotFound            = errors.New("category not found")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type categoryRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (c categoryRepository) Create(ctx context.Context, info *categoryRepoModel.CategoryInfo) (int, error) {
	query := `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id
	`

	var id int
	if err := c.db.QueryRow(ctx, query, &info.Name).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return 0, ErrAlreadyExists
			}
			if pgErr.Code == postgresErrors.InvalidForeignKeyErrCode {
				return 0, ErrInvalidArguments
			}
		}

		return 0, ErrInternalServerError
	}

	return id, nil
}

func (c categoryRepository) GetAll(ctx context.Context, params *model.CategoryGetAllParams) ([]*categoryRepoModel.Category, error) {
	query := `
		SELECT id, name
		FROM categories
		WHERE true
	`

	var args []interface{}
	var conditions []string

	if params.CropId != nil {
		conditions = append(conditions, fmt.Sprintf("crop_id = $%d", len(args)+1))
		args = append(args, *params.CropId)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	if params.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *params.Limit)
	}

	if params.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, *params.Offset)
	}

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}
	defer rows.Close()

	var categories []*categoryRepoModel.Category
	for rows.Next() {
		var cat categoryRepoModel.Category
		if err = rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
		}
		categories = append(categories, &cat)
	}

	return categories, nil
}
