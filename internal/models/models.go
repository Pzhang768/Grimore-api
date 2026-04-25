package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
}

type Team struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time
}

type TeamAgent struct {
	TeamID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	AgentType string         `gorm:"not null"`
	Context   map[string]any `gorm:"serializer:json"`
	Position  int            `gorm:"not null"`
}

type Run struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Status    string    `gorm:"not null"`
	Iteration int       `gorm:"not null;default:0"`
	CreatedAt time.Time
}

type RunEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	RunID     uuid.UUID `gorm:"type:uuid;not null;index"`
	AgentType string    `gorm:"not null"`
	EventType string    `gorm:"not null"`
	Content   string
	CreatedAt time.Time
}

type Deliverable struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	RunID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_deliverables_run_type"`
	Type      string         `gorm:"not null;index:idx_deliverables_run_type"`
	Content   map[string]any `gorm:"serializer:json"`
	CreatedAt time.Time
}
