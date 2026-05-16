package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/internal/auth"
	"server/internal/db"
	"server/internal/lib"
	"server/internal/routers"
)

func main() {
	lib.LoadEnv()
	port := lib.Env.ServerPort

	authInstance, err := auth.NewAuthula()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize auth: %v", err))
	}

	if err := db.InitDB(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	router := routers.Router(authInstance)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		fmt.Printf("Server started on port %d\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Server failed: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	if err := db.DB.Close(); err != nil {
		fmt.Printf("Database connection close error: %v\n", err)
	}

	fmt.Println("Server exited")
}
