package status

import (
	"context"
	"errors"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	"github.com/nogavadu/articles-service/internal/service"
	"github.com/nogavadu/platform_common/pkg/db"
	"log/slog"
)

var (
	ErrAlreadyExists       = errors.New("article already exists")
	ErrInvalidArguments    = errors.New("invalid article arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type statusService struct {
	log *slog.Logger

	statusRepo repository.StatusRepository

	txManager db.TxManager
}

func New(
	log *slog.Logger,
	statusRepo repository.StatusRepository,
	txManager db.TxManager,
) service.StatusService {
	return &statusService{
		log:        log,
		statusRepo: statusRepo,
		txManager:  txManager,
	}
}

func (s *statusService) GetByStatus(ctx context.Context, status string) (*model.Status, error) {
	const op = "statusService.GetByStatus"
	log := s.log.With(slog.String("op", op))

	st, err := s.statusRepo.GetByStatus(ctx, status)
	if err != nil {
		log.Error("failed to get status", "error", err)
		return nil, ErrInternalServerError
	}

	return &model.Status{
		Id:     st.Id,
		Status: st.Status,
	}, nil
}
