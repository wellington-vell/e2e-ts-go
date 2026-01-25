package main

import (
	"fmt"
	"net/http"
	"server/internal"
)

func main() {
	port := internal.Env("SERVER_PORT")

	if err := internal.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	http.HandleFunc("/health", internal.HealthCheck)

	fmt.Printf("Server started on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
