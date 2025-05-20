package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	articleAPI "github.com/nogavadu/articles-service/internal/api/http/article"
	categoryAPI "github.com/nogavadu/articles-service/internal/api/http/category"
	cropAPI "github.com/nogavadu/articles-service/internal/api/http/crop"
	"github.com/nogavadu/articles-service/internal/config/env"
	articleRepo "github.com/nogavadu/articles-service/internal/repository/article"
	categoryRepo "github.com/nogavadu/articles-service/internal/repository/category"
	cropRepo "github.com/nogavadu/articles-service/internal/repository/crop"
	articleServ "github.com/nogavadu/articles-service/internal/service/article"
	categoryServ "github.com/nogavadu/articles-service/internal/service/category"
	cropServ "github.com/nogavadu/articles-service/internal/service/crop"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Error("failed to load pgConfig", slog.String("error", err.Error()))
		os.Exit(1)
	}

	httpServerConfig, err := env.NewHTTPServerConfig()
	if err != nil {
		log.Error("failed to load httpServerConfig", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("configs loaded")

	ctx := context.Background()
	db, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	log.Info("connected to database")

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300, // 5 минут
	}))

	cropRepository := cropRepo.New(db)
	cropService := cropServ.New(log, cropRepository)
	cropApi := cropAPI.New(cropService)

	categoryRepository := categoryRepo.New(db)
	categoryService := categoryServ.New(log, categoryRepository)
	categoryApi := categoryAPI.New(categoryService)

	articleRepository := articleRepo.New(db)
	articleService := articleServ.New(log, articleRepository)
	articleApi := articleAPI.New(articleService)

	router.Route("/api", func(r chi.Router) {
		r.Route("/crops", func(r chi.Router) {
			r.Post("/", cropApi.CreateHandler())
			r.Get("/", cropApi.GetAllHandler())
			r.Get("/{cropId}", cropApi.GetByIdHandler())
			r.Patch("/{cropId}", cropApi.UpdateHandler())
		})

		r.Route("/categories", func(r chi.Router) {
			r.Post("/", categoryApi.CreateHandler())
			r.Get("/", categoryApi.GetAllHandler())
			r.Get("/{categoryId}", categoryApi.GetByIdHandler())
			r.Patch("/{categoryId}", categoryApi.UpdateHandler())
		})

		r.Route("/articles", func(r chi.Router) {
			r.Post("/", articleApi.CreateHandler())
			r.Get("/", articleApi.GetAllHandler())
			r.Get("/{id}", articleApi.GetByIDHandler())
		})
	})

	log.Info("starting server", slog.String("port", strconv.Itoa(httpServerConfig.Port())))
	if err = http.ListenAndServe(httpServerConfig.Address(), router); err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
}
