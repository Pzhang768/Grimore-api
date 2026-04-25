package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestModels_FieldPopulation(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	user := User{ID: id, Email: "test@example.com", CreatedAt: now}
	if user.ID != id || user.Email != "test@example.com" {
		t.Error("User fields not set correctly")
	}

	team := Team{ID: id, UserID: id, Name: "test-team", CreatedAt: now}
	if team.UserID != id || team.Name != "test-team" {
		t.Error("Team fields not set correctly")
	}

	agent := TeamAgent{TeamID: id, AgentType: "fetcher", Position: 1}
	if agent.AgentType != "fetcher" || agent.Position != 1 {
		t.Error("TeamAgent fields not set correctly")
	}

	run := Run{ID: id, TeamID: id, Status: "pending", Iteration: 0, CreatedAt: now}
	if run.Status != "pending" || run.Iteration != 0 {
		t.Error("Run fields not set correctly")
	}

	event := RunEvent{ID: id, RunID: id, AgentType: "fetcher", EventType: "start", Content: "ok", CreatedAt: now}
	if event.EventType != "start" || event.Content != "ok" {
		t.Error("RunEvent fields not set correctly")
	}

	deliverable := Deliverable{ID: id, RunID: id, Type: "resume", CreatedAt: now}
	if deliverable.Type != "resume" {
		t.Error("Deliverable fields not set correctly")
	}
}
