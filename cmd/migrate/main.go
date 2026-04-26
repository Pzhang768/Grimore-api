package main

import (
	"errors"
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

	if err := godotenv.Load(".env.local"); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Warn("failed to load .env.local", "error", err)
	}

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

	// AutoMigrate only adds columns and indexes; destructive changes require manual migration.
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

	fkConstraints := []string{
		`DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_teams_user') THEN ALTER TABLE teams ADD CONSTRAINT fk_teams_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT; END IF; END $$`,
		`DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_team_agents_team') THEN ALTER TABLE team_agents ADD CONSTRAINT fk_team_agents_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE RESTRICT; END IF; END $$`,
		`DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_runs_team') THEN ALTER TABLE runs ADD CONSTRAINT fk_runs_team FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE RESTRICT; END IF; END $$`,
		`DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_run_events_run') THEN ALTER TABLE run_events ADD CONSTRAINT fk_run_events_run FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE RESTRICT; END IF; END $$`,
		`DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_deliverables_run') THEN ALTER TABLE deliverables ADD CONSTRAINT fk_deliverables_run FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE RESTRICT; END IF; END $$`,
	}

	for _, stmt := range fkConstraints {
		if err := db.Exec(stmt).Error; err != nil {
			slog.Error("failed to apply FK constraint", "error", err, "stmt", stmt)
			os.Exit(1)
		}
	}

	slog.Info("migrations complete")
}
