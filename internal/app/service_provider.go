package app

import (
	"context"
	"github.com/nogavadu/articles-service/internal/api/http/article"
	"github.com/nogavadu/articles-service/internal/api/http/category"
	"github.com/nogavadu/articles-service/internal/api/http/crop"
	"github.com/nogavadu/articles-service/internal/client/db"
	"github.com/nogavadu/articles-service/internal/client/db/pg"
	"github.com/nogavadu/articles-service/internal/client/db/transaction"
	"github.com/nogavadu/articles-service/internal/config"
	"github.com/nogavadu/articles-service/internal/config/env"
	"github.com/nogavadu/articles-service/internal/repository"
	articleRepo "github.com/nogavadu/articles-service/internal/repository/article"
	articleImagesRepo "github.com/nogavadu/articles-service/internal/repository/article_images"
	articleRelationsRepo "github.com/nogavadu/articles-service/internal/repository/article_relations"
	categoryRepo "github.com/nogavadu/articles-service/internal/repository/category"
	cropRepo "github.com/nogavadu/articles-service/internal/repository/crop"
	cropCategoriesRepo "github.com/nogavadu/articles-service/internal/repository/crop_categories"
	"github.com/nogavadu/articles-service/internal/service"
	articleServ "github.com/nogavadu/articles-service/internal/service/article"
	categoryServ "github.com/nogavadu/articles-service/internal/service/category"
	cropServ "github.com/nogavadu/articles-service/internal/service/crop"
	"log/slog"
	"os"
)

type serviceProvider struct {
	httpServerConfig config.HTTPServerConfig
	pgConfig         config.PGConfig

	logger *slog.Logger

	cropImpl     *crop.Implementation
	categoryImpl *category.Implementation
	articlesImpl *article.Implementation

	cropService     service.CropService
	categoryService service.CategoryService
	articleService  service.ArticleService

	cropRepository             repository.CropRepository
	categoryRepository         repository.CategoryRepository
	cropsCategoriesRepository  repository.CropCategoriesRepository
	articleRepository          repository.ArticleRepository
	articleImagesRepository    repository.ArticleImagesRepository
	articleRelationsRepository repository.ArticleRelationsRepository

	dbClient  db.Client
	txManager db.TxManager
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (p *serviceProvider) HTTPServerConfig() config.HTTPServerConfig {
	if p.httpServerConfig == nil {
		httpServerConfig, err := env.NewHTTPServerConfig()
		if err != nil {
			p.Logger().Error("failed to get httpServerConfig", slog.String("err", err.Error()))
			panic(err)
		}
		p.httpServerConfig = httpServerConfig
	}
	return p.httpServerConfig
}

func (p *serviceProvider) PGConfig() config.PGConfig {
	if p.pgConfig == nil {
		pgConfig, err := env.NewPGConfig()
		if err != nil {
			p.Logger().Error("failed to get pgConfig", slog.String("err", err.Error()))
			panic(err)
		}
		p.pgConfig = pgConfig
	}
	return p.pgConfig
}

func (p *serviceProvider) Logger() *slog.Logger {
	if p.logger == nil {
		p.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return p.logger
}

func (p *serviceProvider) CropImpl(ctx context.Context) *crop.Implementation {
	if p.cropImpl == nil {
		p.cropImpl = crop.New(p.CropService(ctx))
	}
	return p.cropImpl
}

func (p *serviceProvider) CropService(ctx context.Context) service.CropService {
	if p.cropService == nil {
		p.cropService = cropServ.New(
			p.Logger(),
			p.CropRepository(ctx),
			p.CropCategoriesRepository(ctx),
			p.TxManger(ctx),
		)
	}
	return p.cropService
}

func (p *serviceProvider) CropRepository(ctx context.Context) repository.CropRepository {
	if p.cropRepository == nil {
		p.cropRepository = cropRepo.New(p.DBClient(ctx))
	}
	return p.cropRepository
}

func (p *serviceProvider) CategoryImpl(ctx context.Context) *category.Implementation {
	if p.categoryImpl == nil {
		p.categoryImpl = category.New(p.CategoryService(ctx))
	}
	return p.categoryImpl
}

func (p *serviceProvider) CategoryService(ctx context.Context) service.CategoryService {
	if p.categoryService == nil {
		p.categoryService = categoryServ.New(
			p.Logger(),
			p.CategoryRepository(ctx),
			p.CropCategoriesRepository(ctx),
			p.TxManger(ctx),
		)
	}
	return p.categoryService
}

func (p *serviceProvider) CategoryRepository(ctx context.Context) repository.CategoryRepository {
	if p.categoryRepository == nil {
		p.categoryRepository = categoryRepo.New(p.DBClient(ctx))
	}
	return p.categoryRepository
}

func (p *serviceProvider) CropCategoriesRepository(ctx context.Context) repository.CropCategoriesRepository {
	if p.cropsCategoriesRepository == nil {
		p.cropsCategoriesRepository = cropCategoriesRepo.New(p.DBClient(ctx))
	}
	return p.cropsCategoriesRepository
}

func (p *serviceProvider) ArticleImpl(ctx context.Context) *article.Implementation {
	if p.articlesImpl == nil {
		p.articlesImpl = article.New(p.ArticleService(ctx))
	}
	return p.articlesImpl
}

func (p *serviceProvider) ArticleService(ctx context.Context) service.ArticleService {
	if p.articleService == nil {
		p.articleService = articleServ.New(
			p.Logger(),
			p.ArticleRepository(ctx),
			p.ArticleImagesRepository(ctx),
			p.ArticleRelationsRepository(ctx),
			p.TxManger(ctx),
		)
	}
	return p.articleService
}

func (p *serviceProvider) ArticleRepository(ctx context.Context) repository.ArticleRepository {
	if p.articleRepository == nil {
		p.articleRepository = articleRepo.New(p.DBClient(ctx))
	}
	return p.articleRepository
}

func (p *serviceProvider) ArticleImagesRepository(ctx context.Context) repository.ArticleImagesRepository {
	if p.articleImagesRepository == nil {
		p.articleImagesRepository = articleImagesRepo.New(p.DBClient(ctx))
	}
	return p.articleImagesRepository
}

func (p *serviceProvider) ArticleRelationsRepository(ctx context.Context) repository.ArticleRelationsRepository {
	if p.articleRelationsRepository == nil {
		p.articleRelationsRepository = articleRelationsRepo.New(p.DBClient(ctx))
	}
	return p.articleRelationsRepository
}

func (p *serviceProvider) DBClient(ctx context.Context) db.Client {
	if p.dbClient == nil {
		dbc, err := pg.New(ctx, p.PGConfig().DSN())
		if err != nil {
			p.Logger().Error("failed to create dbClient", slog.String("err", err.Error()))
			panic(err)
		}

		if err = dbc.DB().Ping(ctx); err != nil {
			p.Logger().Error("failed to ping dbClient", slog.String("err", err.Error()))
			panic(err)
		}

		p.dbClient = dbc
	}

	return p.dbClient
}

func (p *serviceProvider) TxManger(ctx context.Context) db.TxManager {
	if p.txManager == nil {
		p.txManager = transaction.NewTransactionManager(p.DBClient(ctx).DB())
	}

	return p.txManager
}
