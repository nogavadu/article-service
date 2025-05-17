package category

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	categoryRepoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
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

func (r *categoryRepository) Create(ctx context.Context, info *categoryRepoModel.CategoryInfo) (int, error) {
	query, args, err := sq.
		Insert("categories").
		PlaceholderFormat(sq.Dollar).
		Columns("name").
		Values(info.Name).
		Suffix("RETURNING id").
		ToSql()

	var id int
	if err = r.db.QueryRow(ctx, query, args...).Scan(&id); err != nil {
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

func (r *categoryRepository) GetAll(ctx context.Context, params *model.CategoryGetAllParams) ([]*categoryRepoModel.Category, error) {
	builder := sq.
		Select("c.id", "c.name").
		PlaceholderFormat(sq.Dollar).
		From("categories AS c")

	if params.CropId != nil {
		builder = builder.Join("article_relations AS ar ON c.id = ar.category_id AND ar.crop_id = ?",
			*params.CropId,
		)
	}

	builder = builder.GroupBy("c.id", "c.name")

	if params.Limit != nil {
		builder = builder.Limit(uint64(*params.Limit))
	}

	if params.Offset != nil {
		builder = builder.Offset(uint64(*params.Offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
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

func (r *categoryRepository) GetById(ctx context.Context, id int) (*categoryRepoModel.Category, error) {
	query, args, err := sq.
		Select("id", "name").
		PlaceholderFormat(sq.Dollar).
		From("categories").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	var cat categoryRepoModel.Category
	if err = pgxscan.Get(ctx, r.db, &cat, query, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}

		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return &cat, nil
}
