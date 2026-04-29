package routers

import (
	"net/http"
	"time"

	"server/internal"
	"server/internal/lib"

	"github.com/Authula/authula"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func Router(auth *authula.Auth) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Timeout(10 * time.Second))

	origins := lib.Env.CorsOrigin
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{origins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With", "Set-Cookie", "Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	r.Get("/api/v1/health", HealthCheck)

	r.Group(func(r chi.Router) {
		r.Get("/spec.json", internal.Docs)
		r.Get("/swagger", internal.SwaggerUI)
		r.Get("/scalar", internal.ScalarUI)
	})

	r.Handle("/auth/*", auth.Handler())

	r.Route("/api/v1/todos", func(r chi.Router) {
		r.Get("/", GetTodos)
		r.Post("/", CreateTodo)
		r.Get("/{id}", GetTodo)
		r.Put("/{id}", UpdateTodo)
		r.Delete("/{id}", DeleteTodo)
	})

	return r
}
