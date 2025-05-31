package status

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/nogavadu/articles-service/internal/repository"
	"github.com/nogavadu/articles-service/internal/repository/status/model"
	"github.com/nogavadu/platform_common/pkg/db"
)

type statusRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repository.StatusRepository {
	return &statusRepository{
		dbc: dbc,
	}
}

func (r *statusRepository) Create(ctx context.Context, status string) (int, error) {
	queryRaw, args, err := sq.
		Insert("entity_status").
		PlaceholderFormat(sq.Dollar).
		Columns("status").
		Values(status).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to make query: %w", err)
	}

	query := db.Query{
		Name:     "statusRepository.Create",
		QueryRaw: queryRaw,
	}

	var id int
	if err = r.dbc.DB().ScanOneContext(ctx, &id, query, args...); err != nil {
		return 0, fmt.Errorf("failed to create status: %w", err)
	}

	return id, nil
}

func (r *statusRepository) GetAll(ctx context.Context) ([]model.Status, error) {
	queryRaw, args, err := sq.
		Select("id", "status").
		PlaceholderFormat(sq.Dollar).
		From("entity_status").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to make query: %w", err)
	}

	query := db.Query{
		Name:     "statusRepository.GetAll",
		QueryRaw: queryRaw,
	}

	var statuses []model.Status
	if err = r.dbc.DB().ScanAllContext(ctx, &statuses, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}

	return statuses, nil
}

func (r *statusRepository) GetByStatus(ctx context.Context, status string) (*model.Status, error) {
	queryRaw, args, err := sq.
		Select("id", "status").
		PlaceholderFormat(sq.Dollar).
		From("entity_status").
		Where(sq.Eq{"status": status}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to make query: %w", err)
	}

	query := db.Query{
		Name:     "statusRepository.GetByStatus",
		QueryRaw: queryRaw,
	}

	var s model.Status
	if err = r.dbc.DB().ScanOneContext(ctx, &s, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	return &s, nil
}

func (r *statusRepository) GetById(ctx context.Context, id int) (*model.Status, error) {
	queryRaw, args, err := sq.
		Select("id", "status").
		PlaceholderFormat(sq.Dollar).
		From("entity_status").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to make query: %w", err)
	}

	query := db.Query{
		Name:     "statusRepository.GetByStatus",
		QueryRaw: queryRaw,
	}

	var s model.Status
	if err = r.dbc.DB().ScanOneContext(ctx, &s, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	return &s, nil
}
