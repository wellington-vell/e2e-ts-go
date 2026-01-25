package internal

import (
	"net/http"

	"github.com/rs/cors"
)

type CORSConf struct {
	AllowedOrigins []string
	MaxAge         int
}

func CORSConfig() *CORSConf {
	return &CORSConf{
		AllowedOrigins: []string{Env("CORS_ORIGIN")},
		MaxAge:         7200, // 2 hours
	}
}

func CORSHandler(handler http.Handler, config *CORSConf) http.Handler {
	if config == nil {
		config = CORSConfig()
	}

	corsOptions := cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           config.MaxAge,
	}

	middleware := cors.New(corsOptions)
	return middleware.Handler(handler)
}
