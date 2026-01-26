package routers

import (
	"maps"

	"github.com/wellington-vell/gorpc"
)

var AllRoutes = func() gorpc.Router {
	routes := make(gorpc.Router)

	maps.Copy(routes, HealthRouter)
	maps.Copy(routes, TodoRouter)

	return routes
}()
