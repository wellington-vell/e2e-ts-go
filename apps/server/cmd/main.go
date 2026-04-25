package main

import (
	"fmt"
	"net/http"

	"server/internal"
	"server/internal/db"
	"server/internal/routers"
)

func main() {
	if err := db.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	port := internal.Env("SERVER_PORT")
	router := routers.Router()

	fmt.Printf("Server started on port %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		panic(fmt.Sprintf("Server failed: %v", err))
	}
}
