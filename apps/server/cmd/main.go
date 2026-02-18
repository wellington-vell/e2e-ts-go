package main

import (
	"fmt"

	"server/internal"
	"server/internal/routers"

	"github.com/wellington-vell/gorpc"
)

func main() {
	port := internal.Env("SERVER_PORT")

	if err := internal.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	app := gorpc.New().
		Prefix("/api").
		Router(routers.AllRoutes).
		Plugin(gorpc.NewScalarPlugin()).
		Plugin(gorpc.NewSwaggerPlugin())

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server starting on port %s\n", port)

	if err := app.ListenAndServe(addr); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
