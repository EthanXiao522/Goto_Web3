package main

import (
	"log/slog"
	"os"

	"github.com/xyd/web3-learning-tracker/internal/config"
	"github.com/xyd/web3-learning-tracker/internal/database"
	"github.com/xyd/web3-learning-tracker/internal/importer"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("usage: go run cmd/seed/main.go <markdown_file>")
		os.Exit(1)
	}
	filePath := os.Args[1]

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

	data, err := importer.Parse(filePath)
	if err != nil {
		slog.Error("failed to parse markdown", "error", err)
		os.Exit(1)
	}

	slog.Info("parsed", "phases", len(data.Phases))
	for _, p := range data.Phases {
		weeks, days, tasks := countPhase(p)
		slog.Info("phase", "num", p.PhaseNumber, "title", p.Title, "weeks", weeks, "days", days, "tasks", tasks)
	}

	if err := importer.Seed(database.DB, data); err != nil {
		slog.Error("failed to seed database", "error", err)
		os.Exit(1)
	}
	slog.Info("seeding complete")
}

func countPhase(p importer.ParsedPhase) (weeks, days, tasks int) {
	weeks = len(p.Weeks)
	for _, w := range p.Weeks {
		days += len(w.Days)
		for _, d := range w.Days {
			tasks += len(d.Tasks)
		}
	}
	return
}
