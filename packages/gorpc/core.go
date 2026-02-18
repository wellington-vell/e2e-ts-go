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
	router         *radixRouter
	prefix         *string
	routers        map[string]Router
	pluginRegistry *pluginRegistry
}

type Router map[string]ProcedureAny

func New() *GORPC {
	return &GORPC{
		router:         NewRouter(),
		routers:        make(map[string]Router),
		pluginRegistry: newPluginRegistry(),
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
	// Create a mux that handles plugin routes first, then falls back to the router
	mux := http.NewServeMux()
	for _, plugin := range g.pluginRegistry.Get() {
		routes := plugin.Routes()
		for path, handler := range routes {
			mux.Handle(path, handler)
		}
	}
	mux.Handle("/", g.router)

	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, mux)
}

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
