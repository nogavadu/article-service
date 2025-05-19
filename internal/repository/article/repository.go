package article

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/lib/postgresErrors"
	"github.com/nogavadu/articles-service/internal/repository"
	articleRepoModel "github.com/nogavadu/articles-service/internal/repository/article/model"
)

var (
	ErrAlreadyExists       = errors.New("article already exists")
	ErrInvalidArguments    = errors.New("invalid arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type articleRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) Create(
	ctx context.Context,
	cropId int,
	categoryId int,
	articleBody *articleRepoModel.ArticleBody,
	images []string,
) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: failed to begin transaction: %w", ErrInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	query, args, err := sq.
		Insert("articles").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "text").
		Values(articleBody.Title, articleBody.Text).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	var articleId int
	if err = tx.QueryRow(ctx, query, args...).Scan(&articleId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return 0, fmt.Errorf("%w: %w", ErrAlreadyExists, err)
			}
		}

		return 0, fmt.Errorf("%w: failed to insert article: %w", ErrInternalServerError, err)
	}

	query, args, err = sq.
		Insert("article_relations").
		PlaceholderFormat(sq.Dollar).
		Columns("crop_id", "category_id", "article_id").
		Values(cropId, categoryId, articleId).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: failed to build relation query: %w", ErrInternalServerError, err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return 0, fmt.Errorf("failed to insert relation: %w: %w", ErrAlreadyExists, err)
			}
			if pgErr.Code == postgresErrors.InvalidForeignKeyErrCode {
				return 0, fmt.Errorf("failed to insert relation: %w: %w", ErrInvalidArguments, err)
			}
		}
		return 0, fmt.Errorf("%w: failed to insert relation: %w", ErrInternalServerError, err)
	}

	err = r.createImages(ctx, articleId, images)
	if err != nil {
		return 0, fmt.Errorf("%w: failed to insert images: %w", ErrInternalServerError, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: failed to commit transaction: %w", ErrInternalServerError, err)
	}

	return articleId, nil
}

func (r *articleRepository) GetById(ctx context.Context, id int) (*articleRepoModel.Article, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to begin transaction: %w", ErrInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	query, args, err := sq.
		Select("id, title, text").
		PlaceholderFormat(sq.Dollar).
		From("articles").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	var article articleRepoModel.Article
	if err = r.db.QueryRow(ctx, query, args...).Scan(&article.Id, &article.Title, &article.Text); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	images, err := r.getImages(ctx, id)
	if err != nil {
		return nil, err
	}
	article.Images = images

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to commit transaction: %w", ErrInternalServerError, err)
	}

	return &article, nil
}

func (r *articleRepository) GetAll(ctx context.Context, params *model.ArticleGetAllParams) ([]*articleRepoModel.Article, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to begin transaction: %w", ErrInternalServerError, err)
	}
	defer tx.Rollback(ctx)

	builder := sq.
		Select("a.id, a.title, a.text").
		PlaceholderFormat(sq.Dollar).
		From("articles AS a")

	if params.CropId != nil && params.CategoryId != nil {
		builder = builder.Join(
			"article_relations AS ar ON a.id = ar.article_id AND ar.crop_id = ? AND ar.category_id = ?",
			*params.CropId, *params.CategoryId,
		)
	} else if params.CropId != nil {
		builder = builder.Join(
			"article_relations AS ar ON a.id = ar.article_id AND ar.crop_id = ?",
			*params.CropId,
		)
	} else if params.CategoryId != nil {
		builder = builder.Join(
			"article_relations AS ar ON a.id = ar.article_id AND ar.category_id = ?",
			*params.CategoryId,
		)
	}

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

	var articles []*articleRepoModel.Article
	for rows.Next() {
		var article articleRepoModel.Article

		if err = rows.Scan(&article.Id, &article.Title, &article.Text); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
		}

		images, err := r.getImages(ctx, article.Id)
		if err != nil {
			return nil, err
		}
		article.Images = images

		articles = append(articles, &article)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to commit transaction: %w", ErrInternalServerError, err)
	}

	return articles, nil
}

func (r *articleRepository) createImages(ctx context.Context, articleId int, images []string) error {
	builder := sq.
		Insert("article_relations").
		PlaceholderFormat(sq.Dollar).
		Columns("article_id", "img")

	for _, image := range images {
		builder = builder.Values(articleId, image)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%w: failed to build query: %w", ErrInternalServerError, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresErrors.AlreadyExistsErrCode {
				return fmt.Errorf("%w: failed to insert image: %w", ErrAlreadyExists, err)
			}
			if pgErr.Code == postgresErrors.InvalidForeignKeyErrCode {
				return fmt.Errorf("%w: failed to insert image: %w", ErrInvalidArguments, err)
			}
		}

		return fmt.Errorf("%w: failed to insert image: %w", ErrInternalServerError, err)
	}

	return nil
}

func (r *articleRepository) getImages(ctx context.Context, articleId int) ([]string, error) {
	query, args, err := sq.
		Select("img").
		PlaceholderFormat(sq.Dollar).
		From("articles_images").
		Where(sq.Eq{"article_id": articleId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
	}
	defer rows.Close()
	var images []string
	for rows.Next() {
		var img string
		if err = rows.Scan(&img); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternalServerError, err)
		}
		images = append(images, img)
	}

	return images, nil
}
