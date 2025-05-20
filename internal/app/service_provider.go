package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nogavadu/articles-service/internal/api/http/article"
	"github.com/nogavadu/articles-service/internal/api/http/category"
	"github.com/nogavadu/articles-service/internal/api/http/crop"
	"github.com/nogavadu/articles-service/internal/config"
	"github.com/nogavadu/articles-service/internal/config/env"
	"github.com/nogavadu/articles-service/internal/repository"
	cropRepo "github.com/nogavadu/articles-service/internal/repository/crop"
	"github.com/nogavadu/articles-service/internal/service"
	cropServ "github.com/nogavadu/articles-service/internal/service/crop"
	"log/slog"
	"os"
)

type serviceProvider struct {
	httpServerConfig config.HTTPServerConfig
	pgConfig         config.PGConfig

	logger *slog.Logger

	cropImpl       *crop.Implementation
	cropService    service.CropService
	cropRepository repository.CropRepository

	categoryImpl       *category.Implementation
	categoryService    service.CategoryService
	categoryRepository repository.CategoryRepository

	articlesImpl      *article.Implementation
	articleService    service.ArticleService
	articleRepository repository.ArticleRepository

	dbClient *pgxpool.Pool
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
		p.cropImpl = crop.New(p.CropService(p.Logger(), ctx))
	}
	return p.cropImpl
}

func (p *serviceProvider) CropService(logger *slog.Logger, ctx context.Context) service.CropService {
	if p.cropService == nil {
		p.cropService = cropServ.New(logger, p.CropRepository(ctx))
	}
	return p.cropService
}

func (p *serviceProvider) CropRepository(ctx context.Context) repository.CropRepository {
	if p.cropRepository == nil {
		p.cropRepository = cropRepo.New(p.DBClient(ctx))
	}
	return p.cropRepository
}

func (p *serviceProvider) DBClient(ctx context.Context) *pgxpool.Pool {
	if p.dbClient == nil {
		dbc, err := pgxpool.New(ctx, p.PGConfig().DSN())
		if err != nil {
			p.Logger().Error("failed to create pgx pool", slog.String("err", err.Error()))
			panic(err)
		}

		if err = dbc.Ping(ctx); err != nil {
			p.Logger().Error("failed to ping db", slog.String("err", err.Error()))
			panic(err)
		}

		p.dbClient = dbc
	}

	return p.dbClient
}
