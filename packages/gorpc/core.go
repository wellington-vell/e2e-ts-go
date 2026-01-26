package gorpc

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type Route struct {
	Method string
	Path   string
}

type Meta struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type GORPC struct {
	router  *radixRouter
	prefix  *string
	routers map[string]Router
}

type Router map[string]ProcedureAny

func New() *GORPC {
	return &GORPC{
		router:  NewRouter(),
		routers: make(map[string]Router),
	}
}

func (g *GORPC) Prefix(path string) *GORPC {
	g.prefix = &path
	return g
}

func (g *GORPC) Router(router Router) *GORPC {
	if router == nil {
		return g
	}

	prefix := "/api"
	if g.prefix != nil {
		prefix = *g.prefix
	}

	for procName, proc := range router {
		// Use reflection to access Route field directly (avoiding getter method)
		procValue := reflect.ValueOf(proc)
		// proc is an interface, so we need to get the underlying value
		if procValue.Kind() == reflect.Interface {
			procValue = procValue.Elem()
		}
		// Now get the actual struct (Procedure is a pointer)
		if procValue.Kind() == reflect.Ptr {
			procValue = procValue.Elem()
		}

		routeField := procValue.FieldByName("Route")
		if !routeField.IsValid() {
			panic(fmt.Sprintf("procedure %s: cannot access route field", procName))
		}

		route, ok := routeField.Interface().(*Route)
		if !ok || route == nil {
			panic(fmt.Sprintf("procedure %s: route is required", procName))
		}

		routePath := route.Path
		if !strings.HasPrefix(routePath, "/") {
			routePath = "/" + routePath
		}
		fullPath := prefix + routePath

		g.router.Insert(fullPath, route.Method, proc)
	}

	return g
}

func (g *GORPC) ListenAndServe(addr string) error {
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, g.router)
}
