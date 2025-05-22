package article_images

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/nogavadu/articles-service/internal/repository"
	"github.com/nogavadu/platform_common/pkg/db"
)

var (
	ErrAlreadyExists       = errors.New("article already exists")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type articleImagesRepository struct {
	dbc db.Client
}

func New(dbc db.Client) repository.ArticleImagesRepository {
	return &articleImagesRepository{
		dbc: dbc,
	}
}

func (r *articleImagesRepository) CreateBulk(ctx context.Context, articleId int, images []string) error {
	builder := sq.
		Insert("articles_images").
		PlaceholderFormat(sq.Dollar).
		Columns("article_id", "img")

	for _, img := range images {
		builder = builder.Values(articleId, img)
	}

	queryRow, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %s: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "articleImagesRepository.CreateBulk",
		QueryRaw: queryRow,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to create article images: %s: %w", ErrInternalServerError, err)
	}

	return nil
}

func (r *articleImagesRepository) GetAll(ctx context.Context, articleId int) ([]string, error) {
	queryRaw, args, err := sq.
		Select("img").
		PlaceholderFormat(sq.Dollar).
		From("articles_images").
		Where(sq.Eq{"article_id": articleId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "articleImagesRepository.GetAll",
		QueryRaw: queryRaw,
	}

	var imgs []string
	if err = r.dbc.DB().ScanAllContext(ctx, &imgs, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get article images: %s: %w", ErrInternalServerError, err)
	}

	return imgs, nil
}

func (r *articleImagesRepository) DeleteBulk(ctx context.Context, articleId int) error {
	queryRaw, args, err := sq.
		Delete("articles_images").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"article_id": articleId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %s: %w", ErrInternalServerError, err)
	}

	query := db.Query{
		Name:     "articleImagesRepository.DeleteBulk",
		QueryRaw: queryRaw,
	}

	if _, err = r.dbc.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete article images: %s: %w", ErrInternalServerError, err)
	}

	return nil
}
