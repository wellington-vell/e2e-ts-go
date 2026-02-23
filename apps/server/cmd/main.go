package main

import (
	"fmt"

	"server/internal"
	"server/internal/routers"

	"github.com/wellington-vell/gorpc"
	"github.com/wellington-vell/gorpc/plugins"
)

func main() {
	port := internal.Env("SERVER_PORT")

	if err := internal.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	app := gorpc.New().
		EnableCORS(gorpc.CORSConfig{
			AllowOrigins:     []string{internal.Env("CORS_ORIGIN")},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
			AllowCredentials: true,
		}).
		Prefix("/api").
		Router(routers.AllRoutes).
		Plugin(plugins.NewScalarPlugin()).
		Plugin(plugins.NewSwaggerPlugin())

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server starting on port %s\n", port)

	if err := app.ListenAndServe(addr); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
