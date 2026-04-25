package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Pzhang768/Grimore-api/config"
	"github.com/Pzhang768/Grimore-api/internal/models"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	_ = godotenv.Load(".env.local")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "error", err)
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("running migrations")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.TeamAgent{},
		&models.Run{},
		&models.RunEvent{},
		&models.Deliverable{},
	); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("migrations complete")
}
