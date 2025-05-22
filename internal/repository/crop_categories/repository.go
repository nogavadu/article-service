package crop_categories

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/nogavadu/articles-service/internal/client/db"
	"github.com/nogavadu/articles-service/internal/repository"
)

var (
	ErrAlreadyExists       = errors.New("crop already exists")
	ErrNotFound            = errors.New("crop not found")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type cropCategoriesRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repository.CropCategoriesRepository {
	return &cropCategoriesRepository{
		dbc: dbc,
	}
}

func (r *cropCategoriesRepository) Create(ctx context.Context, cropId int, categoriesId int) error {
	queryRaw, args, err := sq.
		Insert("crops_categories").
		PlaceholderFormat(sq.Dollar).
		Columns("crop_id", "category_id").
		Values(cropId, categoriesId).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "cropCategories.create",
		QueryRaw: queryRaw,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", ErrInternalServerError, err)
	}

	return nil
}

func (r *cropCategoriesRepository) Delete(ctx context.Context, cropId int, categoryId int) error {
	queryRaw, args, err := sq.
		Delete("crops_categories").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			"crop_id":     cropId,
			"category_id": categoryId,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "cropCategories.delete",
		QueryRaw: queryRaw,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s: %w", ErrInternalServerError, err)
	}

	return nil
}
