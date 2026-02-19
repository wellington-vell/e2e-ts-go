package gorpc

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Route struct {
	Method string
	Path   string
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

type Meta struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

// GORPC is the main framework instance that manages routers, plugins, and the HTTP server.
// It provides a fluent API for registering procedures, mounting plugins, and starting the server.
type GORPC struct {
	router         *radixRouter
	prefix         *string
	routers        map[string]Router
	pluginRegistry *pluginRegistry
	corsConfig     *CORSConfig
}

type Router map[string]ProcedureAny

// New creates a new GORPC instance with an empty router and plugin registry.
// This is the entry point for creating a new gorpc application.
func New() *GORPC {
	return &GORPC{
		router:         NewRouter(),
		routers:        make(map[string]Router),
		pluginRegistry: newPluginRegistry(),
	}
}

func (g *GORPC) EnableCORS(config CORSConfig) *GORPC {
	g.corsConfig = &config
	return g
}

func (g *GORPC) Prefix(path string) *GORPC {
	g.prefix = &path
	return g
}

// Router registers a collection of procedures under a common prefix.
// It uses reflection to extract route metadata from each procedure struct,
// including the HTTP method, path, and path parameters. This design allows
// procedures to define their routing behavior declaratively while keeping
// the registration process centralized.
//
// The prefix defaults to "/api" if not explicitly set via Prefix().
// Reflection is used here to avoid requiring getter methods on procedure structs,
// following the convention of direct field access for configuration.
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

		// Extract path parameters from the route path
		pathParams := extractPathParamsFromRoute(routePath)
		pathParamsField := procValue.FieldByName("PathParams")
		if pathParamsField.IsValid() && pathParamsField.Kind() == reflect.Slice {
			pathParamsField.Set(reflect.ValueOf(pathParams))
		}

		g.router.Insert(fullPath, route.Method, proc)
	}

	return g
}

// Plugin registers a plugin with the gorpc instance. Plugins can extend
// functionality by registering HTTP routes, adding middleware, or modifying
// the application behavior. The plugin's Register method is called immediately
// to allow it to access the GORPC instance and register its routes.
//
// Panics if plugin registration fails, as this indicates a configuration
// error that should be caught during startup rather than at runtime.
func (g *GORPC) Plugin(plugin Plugin) *GORPC {
	if plugin == nil {
		return g
	}
	g.pluginRegistry.Register(plugin)
	if err := plugin.Register(g); err != nil {
		panic(fmt.Sprintf("Failed to register plugin %s: %v", plugin.Name(), err))
	}
	return g
}

func (g *GORPC) hasPlugin(name string) bool {
	for _, plugin := range g.pluginRegistry.Get() {
		if plugin.Name() == name {
			return true
		}
	}
	return false
}

type ProcedureInfo struct {
	Path       string
	Method     string
	Route      *Route
	Meta       Meta
	Tags       []string
	InputType  reflect.Type
	OutputType reflect.Type
	ErrorCodes []int
	PathParams []string
}

func (g *GORPC) GetRouters() map[string]Router {
	return g.routers
}

// GetAllProcedures traverses the radix router and returns all registered procedures with their metadata
func (g *GORPC) GetAllProcedures() []ProcedureInfo {
	var procedures []ProcedureInfo
	routerProcs := g.router.GetAllProcedures()

	for _, routerProc := range routerProcs {
		info := g.extractProcedureInfo(routerProc.Procedure, routerProc.Path, routerProc.Method)
		if info != nil {
			procedures = append(procedures, *info)
		}
	}

	return procedures
}

// extractProcedureInfo extracts metadata from a procedure using reflection
func (g *GORPC) extractProcedureInfo(proc ProcedureAny, path, method string) *ProcedureInfo {
	procValue := reflect.ValueOf(proc)
	if procValue.Kind() == reflect.Interface {
		procValue = procValue.Elem()
	}
	if procValue.Kind() == reflect.Ptr {
		procValue = procValue.Elem()
	}

	if !procValue.IsValid() {
		return nil
	}

	info := &ProcedureInfo{
		Path:   path,
		Method: method,
	}

	routeField := procValue.FieldByName("Route")
	if routeField.IsValid() {
		if route, ok := routeField.Interface().(*Route); ok {
			info.Route = route
		}
	}

	metaField := procValue.FieldByName("Meta")
	if metaField.IsValid() {
		if meta, ok := metaField.Interface().(Meta); ok {
			info.Meta = meta
		}
	}

	tagsField := procValue.FieldByName("Tags")
	if tagsField.IsValid() {
		if tags, ok := tagsField.Interface().([]string); ok {
			info.Tags = tags
		}
	}

	inputTypeField := procValue.FieldByName("InputType")
	if inputTypeField.IsValid() {
		if inputType, ok := inputTypeField.Interface().(reflect.Type); ok {
			info.InputType = inputType
		}
	}

	outputTypeField := procValue.FieldByName("OutputType")
	if outputTypeField.IsValid() {
		if outputType, ok := outputTypeField.Interface().(reflect.Type); ok {
			info.OutputType = outputType
		}
	}

	errorCodesField := procValue.FieldByName("ErrorCodes")
	if errorCodesField.IsValid() {
		if errorCodes, ok := errorCodesField.Interface().([]int); ok {
			info.ErrorCodes = errorCodes
		}
	}

	pathParamsField := procValue.FieldByName("PathParams")
	if pathParamsField.IsValid() {
		if pathParams, ok := pathParamsField.Interface().([]string); ok {
			info.PathParams = pathParams
		}
	}

	return info
}

func (g *GORPC) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	for _, plugin := range g.pluginRegistry.Get() {
		routes := plugin.Routes()
		for path, handler := range routes {
			mux.Handle(path, g.wrapCORS(handler))
		}
	}
	mux.Handle("/", g.wrapCORS(g.router))

	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, mux)
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
	hasWildcard := false
	for _, o := range cfg.AllowOrigins {
		if o == "*" {
			hasWildcard = true
			break
		}
	}

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

// extractPathParamsFromRoute extracts parameter names from route paths.
// It identifies path parameters by the colon prefix (e.g., :id in /todos/:id).
// This is used to pre-populate the PathParams field on procedures so that
// parameter names are available for documentation and validation without
// requiring runtime extraction for every request.
func extractPathParamsFromRoute(path string) []string {
	var params []string
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		if len(segment) > 0 && segment[0] == ':' {
			params = append(params, segment[1:])
		}
	}
	return params
}
