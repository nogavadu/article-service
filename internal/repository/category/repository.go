package category

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	categoryRepoModel "github.com/nogavadu/articles-service/internal/repository/category/model"
	"github.com/nogavadu/platform_common/pkg/db"
	"time"
)

var (
	ErrAlreadyExists       = errors.New("category already exists")
	ErrNotFound            = errors.New("category not found")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type categoryRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repository.CategoryRepository {
	return &categoryRepository{
		dbc: dbc,
	}
}

func (r *categoryRepository) Create(ctx context.Context, info *categoryRepoModel.CategoryInfo) (int, error) {
	queryRaw, args, err := sq.
		Insert("categories").
		PlaceholderFormat(sq.Dollar).
		Columns(
			"name",
			"description",
			"icon",
			"author",
			"status",
			"created_at",
			"updated_at",
		).
		Values(
			info.Name,
			info.Description,
			info.Icon,
			info.Author,
			info.Status,
			time.Now(),
			time.Now(),
		).
		Suffix("RETURNING id").
		ToSql()

	query := db.Query{
		Name:     "categoryRepository.Create",
		QueryRaw: queryRaw,
	}

	var id int
	if err = r.dbc.DB().ScanOneContext(ctx, &id, query, args...); err != nil {
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

func (r *categoryRepository) GetAll(
	ctx context.Context,
	params *categoryRepoModel.CategoryGetAllParams,
) ([]categoryRepoModel.Category, error) {
	builder := sq.
		Select(
			"c.id",
			"c.name",
			"c.description",
			"c.icon",
			"c.author",
			"c.status",
			"c.created_at",
			"c.updated_at",
		).
		PlaceholderFormat(sq.Dollar).
		From("categories AS c")

	if params.CropId != nil {
		builder = builder.Join("crops_categories AS ar ON c.id = ar.category_id AND ar.crop_id = ?",
			*params.CropId,
		)
	}

	builder = builder.
		Where(sq.Eq{"c.status": params.Status}).
		GroupBy("c.id", "c.name")

	queryRaw, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "categoryRepository.GetAll",
		QueryRaw: queryRaw,
	}

	var categories []categoryRepoModel.Category
	if err = r.dbc.DB().ScanAllContext(ctx, &categories, query, args...); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return categories, nil
}

func (r *categoryRepository) GetById(ctx context.Context, id int) (*categoryRepoModel.Category, error) {
	queryRaw, args, err := sq.
		Select(
			"id",
			"name",
			"description",
			"icon",
			"author",
			"status",
			"created_at",
			"updated_at",
		).
		PlaceholderFormat(sq.Dollar).
		From("categories").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "categoryRepository.GetById",
		QueryRaw: queryRaw,
	}

	var category categoryRepoModel.Category
	if err = r.dbc.DB().ScanOneContext(ctx, &category, query, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}

		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, id int, input *categoryRepoModel.UpdateInput) error {
	values := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.Name != nil {
		values["name"] = *input.Name
	}
	if input.Description != nil {
		values["description"] = *input.Description
	}
	if input.Icon != nil {
		values["icon"] = *input.Icon
	}
	if input.Status != nil {
		values["status"] = *input.Status
	}

	queryRaw, args, err := sq.
		Update("categories").
		PlaceholderFormat(sq.Dollar).
		SetMap(values).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "categoryRepository.Update",
		QueryRaw: queryRaw,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	queryRaw, args, err := sq.
		Delete("categories").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "categoryRepository.Delete",
		QueryRaw: queryRaw,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return nil
}
