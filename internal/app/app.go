package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	articleAPI "github.com/nogavadu/articles-service/internal/api/http/article"
	categoryAPI "github.com/nogavadu/articles-service/internal/api/http/category"
	cropAPI "github.com/nogavadu/articles-service/internal/api/http/crop"
	articleRepo "github.com/nogavadu/articles-service/internal/repository/article"
	categoryRepo "github.com/nogavadu/articles-service/internal/repository/category"
	cropRepo "github.com/nogavadu/articles-service/internal/repository/crop"
	articleServ "github.com/nogavadu/articles-service/internal/service/article"
	categoryServ "github.com/nogavadu/articles-service/internal/service/category"
	cropServ "github.com/nogavadu/articles-service/internal/service/crop"
	"log/slog"
	"net/http"
	"strconv"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *chi.Mux
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return a.runHttpServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initHttpServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300, // 5 минут
	}))

	router.Route("/api", func(r chi.Router) {
		a.initCropAPI(ctx, r)
		a.initCategoryAPI(ctx, r)
		a.initArticleAPI(ctx, r)
	})

	a.httpServer = router

	return nil
}

func (a *App) initCropAPI(ctx context.Context, r chi.Router) {
	cropRepository := cropRepo.New(a.serviceProvider.DBClient(ctx))
	cropService := cropServ.New(a.serviceProvider.Logger(), cropRepository)
	cropApi := cropAPI.New(cropService)

	r.Route("/crops", func(r chi.Router) {
		r.Post("/", cropApi.CreateHandler())
		r.Get("/", cropApi.GetAllHandler())
		r.Get("/{cropId}", cropApi.GetByIdHandler())
		r.Patch("/{cropId}", cropApi.UpdateHandler())
	})
}

func (a *App) initCategoryAPI(ctx context.Context, r chi.Router) {
	categoryRepository := categoryRepo.New(a.serviceProvider.DBClient(ctx))
	categoryService := categoryServ.New(a.serviceProvider.Logger(), categoryRepository)
	categoryApi := categoryAPI.New(categoryService)

	r.Route("/categories", func(r chi.Router) {
		r.Post("/", categoryApi.CreateHandler())
		r.Get("/", categoryApi.GetAllHandler())
		r.Get("/{categoryId}", categoryApi.GetByIdHandler())
		r.Patch("/{categoryId}", categoryApi.UpdateHandler())
	})
}

func (a *App) initArticleAPI(ctx context.Context, r chi.Router) {
	articleRepository := articleRepo.New(a.serviceProvider.DBClient(ctx))
	articleService := articleServ.New(a.serviceProvider.Logger(), articleRepository)
	articleApi := articleAPI.New(articleService)

	r.Route("/articles", func(r chi.Router) {
		r.Post("/", articleApi.CreateHandler())
		r.Get("/", articleApi.GetAllHandler())
		r.Get("/{id}", articleApi.GetByIDHandler())
	})
}

func (a *App) runHttpServer() error {
	a.serviceProvider.Logger().Info(
		"starting server", slog.String("port", strconv.Itoa(a.serviceProvider.HTTPServerConfig().Port())),
	)

	if err := http.ListenAndServe(
		a.serviceProvider.HTTPServerConfig().Address(), a.httpServer,
	); err != nil {
		a.serviceProvider.Logger().Error("failed to start server", slog.String("error", err.Error()))
	}

	return nil
}
