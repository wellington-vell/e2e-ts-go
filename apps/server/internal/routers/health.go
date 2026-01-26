package routers

import (
	"github.com/wellington-vell/gorpc"
)

var HealthRouter = gorpc.Router{
	"healthCheck": gorpc.OS().
		Output("").
		Tag("health").
		Meta(gorpc.Meta{
			Summary:     "Health check endpoint",
			Description: "Returns OK if the server is running",
		}).
		Route(gorpc.Route{
			Method: "GET",
			Path:   "/health",
		}).
		Handler(func(ctx *gorpc.Context, input interface{}) (interface{}, error) {
			return "OK", nil
		}).
		Build(),
}
