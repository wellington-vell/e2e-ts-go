package routers

import (
	"net/http"

	"server/internal"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{internal.Env("CORS_ORIGIN")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	r.Get("/api/v1/health", HealthCheck)

	r.Group(func(r chi.Router) {
		r.Get("/spec.json", internal.Docs)
		r.Get("/swagger", internal.SwaggerUI)
		r.Get("/scalar", internal.ScalarUI)
	})

	r.Route("/api/v1/todos", func(r chi.Router) {
		r.Get("/", HandleGetTodos)
		r.Post("/", HandleCreateTodo)
		r.Get("/{id}", HandleGetTodo)
		r.Put("/{id}", HandleUpdateTodo)
		r.Delete("/{id}", HandleDeleteTodo)
	})

	return r
}
