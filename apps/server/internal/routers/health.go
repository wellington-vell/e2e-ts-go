package routers

import (
	"github.com/wellington-vell/gorpc"
)

var HealthRouter = gorpc.Router{
	"healthCheck": gorpc.OS[struct{}, string]().
		Tag("health").
		Meta(gorpc.Meta{
			Summary:     "Health check endpoint",
			Description: "Returns OK if the server is running",
		}).
		Route(gorpc.Route{
			Method: "GET",
			Path:   "/health",
		}).
		Handler(func(ctx *gorpc.Context, input struct{}) (string, error) {
			return "OK", nil
		}).
		Build(),
}
