package main

import (
	"log/slog"
	"os"

	"expenses/internal/config"
	"expenses/internal/database"
	"expenses/internal/server"
)

func main() {
	cfg := config.Load()

	if err := database.RunMigrations(cfg.MigrationsPath, cfg.DBURL); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	db, err := database.NewPostgreSQL(cfg.DBURL)
	if err != nil {
		slog.Error("failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	e := server.New(cfg, db)

	slog.Info("server starting", "address", cfg.ServerAddress)
	if err := e.Start(cfg.ServerAddress); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
