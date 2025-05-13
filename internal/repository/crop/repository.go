package crop

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	cropRepoModel "github.com/nogavadu/articles-service/internal/repository/crop/model"
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

func (r *cropRepository) Create(ctx context.Context, info *cropRepoModel.CropInfo) (int, error) {
	const op = "cropRepository.Create"

	query := `
		INSERT INTO crops (name)
		VALUES ($1)
		RETURNING id
   	`

	var cropId int
	if err := r.db.QueryRow(ctx, query, info.Name).Scan(&cropId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return 0, fmt.Errorf("%s: %w: %w", op, ErrAlreadyExists, err)
			}
		}

		return 0, fmt.Errorf("%s: %w: %w", op, ErrInternalServerError, err)
	}

	return cropId, nil
}

func (r *cropRepository) GetAll(ctx context.Context) ([]*cropRepoModel.Crop, error) {
	const op = "cropRepository.GetAll"

	query := `
		SELECT id, name FROM crops
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrInternalServerError, err)
	}
	defer rows.Close()

	var crops []*cropRepoModel.Crop
	for rows.Next() {
		var crop cropRepoModel.Crop
		var cropInfo cropRepoModel.CropInfo

		if err = rows.Scan(&crop.ID, &cropInfo.Name); err != nil {
			return nil, fmt.Errorf("%s: %w: %w", op, ErrInternalServerError, err)
		}

		crop.Info = &cropInfo

		crops = append(crops, &crop)
	}

	return crops, nil
}
