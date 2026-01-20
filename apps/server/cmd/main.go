package main

import (
	"fmt"
	"net/http"
	"server/internal"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK!")
}

func main() {
	port := internal.Env("SERVER_PORT")

	http.HandleFunc("/health", healthHandler)

	fmt.Printf("Server started on port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
