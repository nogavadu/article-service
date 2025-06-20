package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/nogavadu/articles-service/internal/middlewares"
	"github.com/nogavadu/platform_common/pkg/closer"
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
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

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

func (a *App) initAuthAPI(r chi.Router) {
	authApi := a.serviceProvider.AuthImpl()

	r.Post("/register", authApi.RegisterHandler())
	r.Post("/login", authApi.LoginHandler())

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Get("/refreshToken", authApi.GetRefreshTokenHandler())
	})
}

func (a *App) initUserAPI(r chi.Router) {
	userApi := a.serviceProvider.UserImpl()

	r.Route("/users", func(r chi.Router) {
		r.Get("/{userId}", userApi.GetByIdHandler())

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)

			r.Patch("/{userId}", userApi.UpdateHandler())
		})
	})
}

func (a *App) initCropAPI(ctx context.Context, r chi.Router) {
	cropApi := a.serviceProvider.CropImpl(ctx)

	r.Route("/crops", func(r chi.Router) {
		r.Get("/", cropApi.GetAllHandler())
		r.Get("/{cropId}", cropApi.GetByIdHandler())

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)

			r.Post("/", cropApi.CreateHandler())
			r.Patch("/{cropId}", cropApi.UpdateHandler())
			r.Delete("/{cropId}", cropApi.DeleteHandler())

			r.Post("/{cropId}/{categoryId}", cropApi.AddRelationHandler())
			r.Delete("/{cropId}/{categoryId}", cropApi.RemoveRelationHandler())
		})
	})
}

func (a *App) initCategoryAPI(ctx context.Context, r chi.Router) {
	categoryApi := a.serviceProvider.CategoryImpl(ctx)

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", categoryApi.GetAllHandler())
		r.Get("/{categoryId}", categoryApi.GetByIdHandler())

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)

			r.Post("/", categoryApi.CreateHandler())
			r.Patch("/{categoryId}", categoryApi.UpdateHandler())
			r.Delete("/{categoryId}", categoryApi.DeleteHandler())
		})
	})
}

func (a *App) initArticleAPI(ctx context.Context, r chi.Router) {
	articleApi := a.serviceProvider.ArticleImpl(ctx)

	r.Route("/articles", func(r chi.Router) {
		r.Get("/", articleApi.GetAllHandler())
		r.Get("/{articleId}", articleApi.GetByIDHandler())

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)

			r.Post("/", articleApi.CreateHandler())
			r.Patch("/{articleId}", articleApi.UpdateHandler())
			r.Delete("/{articleId}", articleApi.DeleteHandler())
		})
	})
}

func (a *App) initHttpServer(ctx context.Context) error {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300, // 5 минут
	}))

	router.Route("/api", func(r chi.Router) {
		a.initAuthAPI(r)
		a.initUserAPI(r)
		a.initCropAPI(ctx, r)
		a.initCategoryAPI(ctx, r)
		a.initArticleAPI(ctx, r)
	})

	a.httpServer = router

	return nil
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
