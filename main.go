package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/brix101/linkr/internal/api"
	"github.com/brix101/linkr/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found, proceeding with environment variables")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := db.Connect(ctx, 5)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}
	defer pool.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	api := api.New(ctx, pool)
	srv := api.Server(port)

	go func() { _ = srv.ListenAndServe() }()

	slog.Info("🚀🚀🚀 Server started", slog.String("port", port))
	<-ctx.Done()

	_ = srv.Shutdown(ctx)
}
