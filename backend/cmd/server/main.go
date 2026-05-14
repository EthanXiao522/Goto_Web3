package main

import (
	"log/slog"
	"os"

	"github.com/xyd/web3-learning-tracker/internal/config"
	"github.com/xyd/web3-learning-tracker/internal/database"
	"github.com/xyd/web3-learning-tracker/internal/router"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg.DBDSN); err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		slog.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	r := router.Setup(database.DB, cfg.JWTSecret)
	slog.Info("server starting", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
