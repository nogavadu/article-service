package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	articleAPI "github.com/nogavadu/articles-service/internal/api/http/article"
	cropAPI "github.com/nogavadu/articles-service/internal/api/http/crop"
	"github.com/nogavadu/articles-service/internal/config/env"
	articleRepo "github.com/nogavadu/articles-service/internal/repository/article"
	cropRepo "github.com/nogavadu/articles-service/internal/repository/crop"
	articleServ "github.com/nogavadu/articles-service/internal/service/article"
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

	cropRepository := cropRepo.New(db)
	cropService := cropServ.New(log, cropRepository)
	cropApi := cropAPI.New(cropService)

	articleRepository := articleRepo.New(db)
	articleService := articleServ.New(log, articleRepository)
	articleApi := articleAPI.New(articleService)

	router.Route("/api", func(r chi.Router) {
		r.Route("/{crop_id}", func(r chi.Router) {
			r.Route("/{category_id}", func(r chi.Router) {
				r.Get("/", articleApi.GetListHandler())
				r.Get("/{article_id}", articleApi.GetByIDHandler())
			})
		})

		r.Route("/crops", func(r chi.Router) {
			r.Post("/", cropApi.CreateHandler())
			r.Get("/", cropApi.GetAllHandler())
		})

		r.Route("/articles", func(r chi.Router) {
			r.Post("/", articleApi.CreateHandler())
		})
	})

	log.Info("starting server", slog.String("port", strconv.Itoa(httpServerConfig.Port())))
	if err = http.ListenAndServe(httpServerConfig.Address(), router); err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
}
