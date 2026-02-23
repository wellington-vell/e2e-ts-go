package gorpc

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

func (g *GORPC) EnableCORS(config CORSConfig) *GORPC {
	g.corsConfig = &config
	return g
}

func (g *GORPC) wrapCORS(handler http.Handler) http.Handler {
	if g.corsConfig == nil {
		return handler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if r.Method == "OPTIONS" {
			g.handleCORS(w, r, origin)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		g.handleCORS(w, r, origin)
		handler.ServeHTTP(w, r)
	})
}

func (g *GORPC) handleCORS(w http.ResponseWriter, _ *http.Request, origin string) {
	cfg := g.corsConfig

	allowOrigin := ""
	hasWildcard := slices.Contains(cfg.AllowOrigins, "*")

	if hasWildcard {
		allowOrigin = "*"
	} else if origin != "" {
		for _, o := range cfg.AllowOrigins {
			if o == origin {
				allowOrigin = o
				break
			}
		}
	}

	if allowOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
	}

	if len(cfg.AllowMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
	}

	if len(cfg.AllowHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
	}

	if cfg.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if len(cfg.ExposeHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
	}

	if cfg.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
	}
}
