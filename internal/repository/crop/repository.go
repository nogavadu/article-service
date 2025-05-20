package crop

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	cropRepoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
	"time"
)

var (
	ErrAlreadyExists       = errors.New("crop already exists")
	ErrNotFound            = errors.New("crop not found")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type cropRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.CropRepository {
	return &cropRepository{
		db: db,
	}
}

func (r *cropRepository) Create(ctx context.Context, cropInfo *cropRepoModel.CropInfo) (int, error) {
	query, args, err := sq.
		Insert("crops").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "description", "img", "created_at", "updated_at").
		Values(cropInfo.Name, cropInfo.Description, cropInfo.Img, time.Now(), time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	fmt.Printf("REPO CROP INFO: %s\n", cropInfo)
	fmt.Print(query)

	var cropId int
	if err = r.db.QueryRow(ctx, query, args...).Scan(&cropId); err != nil {
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

func (r *cropRepository) GetAll(ctx context.Context) ([]*cropRepoModel.Crop, error) {
	query, _, err := sq.
		Select("id", "name", "description", "img", "created_at", "updated_at").
		From("crops").
		ToSql()

	var crops []*cropRepoModel.Crop
	if err = pgxscan.Select(ctx, r.db, &crops, query); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return crops, nil
}

func (r *cropRepository) GetById(ctx context.Context, id int) (*cropRepoModel.Crop, error) {
	query, args, err := sq.
		Select("id", "name", "description", "img", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("crops").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	var crop cropRepoModel.Crop
	if err = pgxscan.Get(ctx, r.db, &crop, query, args...); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return &crop, nil
}

func (r *cropRepository) Update(ctx context.Context, id int, input *cropRepoModel.UpdateInput) error {
	builder := sq.Update("crops").
		PlaceholderFormat(sq.Dollar)

	if input.Name != nil {
		builder = builder.Set("name", *input.Name)
	}
	if input.Description != nil {
		builder = builder.Set("description", *input.Description)
	}
	if input.Img != nil {
		builder = builder.Set("img", *input.Img)
	}

	query, args, err := builder.
		Set("updated_ad", time.Now()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	return nil
}
