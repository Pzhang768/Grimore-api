package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type Team struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index;foreignKey:UserID;constraint:OnDelete:RESTRICT"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time
}

func (t *Team) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TeamAgent struct {
	TeamID    uuid.UUID      `gorm:"type:uuid;primaryKey;foreignKey:TeamID;constraint:OnDelete:RESTRICT"`
	AgentType string         `gorm:"primaryKey"`
	Context   map[string]any `gorm:"serializer:json"`
	Position  int            `gorm:"not null"`
}

type Run struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null;index;foreignKey:TeamID;constraint:OnDelete:RESTRICT"`
	Status    string    `gorm:"not null"`
	Iteration int       `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Run) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type RunEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	RunID     uuid.UUID `gorm:"type:uuid;not null;index;foreignKey:RunID;constraint:OnDelete:RESTRICT"`
	AgentType string    `gorm:"not null"`
	EventType string    `gorm:"not null"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time
}

func (re *RunEvent) BeforeCreate(_ *gorm.DB) error {
	if re.ID == uuid.Nil {
		re.ID = uuid.New()
	}
	return nil
}

type Deliverable struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	RunID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_deliverables_run_type;foreignKey:RunID;constraint:OnDelete:RESTRICT"`
	Type      string         `gorm:"not null;index:idx_deliverables_run_type"`
	Content   map[string]any `gorm:"serializer:json;not null;default:'{}'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Deliverable) BeforeCreate(_ *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
