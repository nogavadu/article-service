package crop

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	cropRepoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
	"github.com/nogavadu/platform_common/pkg/db"
	"time"
)

var (
	ErrAlreadyExists       = errors.New("crop already exists")
	ErrNotFound            = errors.New("crop not found")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type cropRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repository.CropRepository {
	return &cropRepository{
		dbc: dbc,
	}
}

func (r *cropRepository) Create(ctx context.Context, cropInfo *cropRepoModel.CropInfo) (int, error) {
	queryRaw, args, err := sq.
		Insert("crops").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "description", "img", "created_at", "updated_at").
		Values(cropInfo.Name, cropInfo.Description, cropInfo.Img, time.Now(), time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "cropRepository.Create",
		QueryRaw: queryRaw,
	}

	var cropId int
	if err = r.dbc.DB().ScanOneContext(ctx, &cropId, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return 0, fmt.Errorf("%w: %w", ErrAlreadyExists, err)
			}
		}

		return 0, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return cropId, nil
}

func (r *cropRepository) GetAll(ctx context.Context) ([]cropRepoModel.Crop, error) {
	queryRaw, _, err := sq.
		Select("id", "name", "description", "img", "created_at", "updated_at").
		From("crops").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "cropRepository.GetAll",
		QueryRaw: queryRaw,
	}

	var crops []cropRepoModel.Crop
	if err = r.dbc.DB().ScanAllContext(ctx, &crops, query); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return crops, nil
}

func (r *cropRepository) GetById(ctx context.Context, id int) (*cropRepoModel.Crop, error) {
	queryRaw, args, err := sq.
		Select("id", "name", "description", "img", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("crops").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "cropRepository.GetById",
		QueryRaw: queryRaw,
	}

	var crop cropRepoModel.Crop
	if err = r.dbc.DB().ScanOneContext(ctx, &crop, query, args...); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return &crop, nil
}

func (r *cropRepository) Update(ctx context.Context, id int, input *cropRepoModel.UpdateInput) error {
	builder := sq.Update("crops").PlaceholderFormat(sq.Dollar)

	if input.Name != nil {
		builder = builder.Set("name", *input.Name)
	}
	if input.Description != nil {
		builder = builder.Set("description", *input.Description)
	}
	if input.Img != nil {
		builder = builder.Set("img", *input.Img)
	}

	queryRaw, args, err := builder.Set("updated_at", time.Now()).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "cropRepository.Update",
		QueryRaw: queryRaw,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return nil
}

func (r *cropRepository) Delete(ctx context.Context, id int) error {
	queryRaw, args, err := sq.
		Delete("crops").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "cropRepository.Delete",
		QueryRaw: queryRaw,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return nil
}
