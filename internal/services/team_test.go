package services

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}
	return db, mock
}

func TestListTeams_Success(t *testing.T) {
	db, mock := newMockDB(t)
	svc := NewTeamService(db)
	userID := uuid.New()
	teamID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"}).
		AddRow(teamID, userID, "team one", now).
		AddRow(uuid.New(), userID, "team two", now)

	mock.ExpectQuery(`SELECT \* FROM "teams" WHERE user_id`).
		WithArgs(userID).
		WillReturnRows(rows)

	teams, err := svc.ListTeams(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(teams) != 2 {
		t.Errorf("expected 2 teams, got %d", len(teams))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestListTeams_Empty(t *testing.T) {
	db, mock := newMockDB(t)
	svc := NewTeamService(db)
	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at"})
	mock.ExpectQuery(`SELECT \* FROM "teams" WHERE user_id`).
		WithArgs(userID).
		WillReturnRows(rows)

	teams, err := svc.ListTeams(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if teams == nil || len(teams) != 0 {
		t.Errorf("expected empty slice, got %d", len(teams))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestListTeams_DBError(t *testing.T) {
	db, mock := newMockDB(t)
	svc := NewTeamService(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "teams" WHERE user_id`).
		WithArgs(userID).
		WillReturnError(errors.New("db error"))

	_, err := svc.ListTeams(context.Background(), userID)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

var baseAgents = []CreateAgentInput{
	{AgentType: "fetcher", Position: 0},
}

func TestCreateTeam_Validation(t *testing.T) {
	tests := []struct {
		name     string
		teamName string
		agents   []CreateAgentInput
		wantErr  error
	}{
		{
			name:     "empty name",
			teamName: "",
			agents:   baseAgents,
			wantErr:  ErrTeamNameRequired,
		},
		{
			name:     "no agents",
			teamName: "my team",
			agents:   []CreateAgentInput{},
			wantErr:  ErrAgentsRequired,
		},
		{
			name:     "too many agents",
			teamName: "my team",
			agents: []CreateAgentInput{
				{AgentType: "fetcher", Position: 0},
				{AgentType: "analyser", Position: 1},
				{AgentType: "tailor", Position: 2},
				{AgentType: "coordinator", Position: 3},
				{AgentType: "fetcher", Position: 4},
			},
			wantErr: ErrTooManyAgents,
		},
		{
			name:     "invalid agent type",
			teamName: "my team",
			agents:   []CreateAgentInput{{AgentType: "unknown", Position: 0}},
			wantErr:  ErrInvalidAgentType,
		},
		{
			name:     "duplicate position",
			teamName: "my team",
			agents: []CreateAgentInput{
				{AgentType: "fetcher", Position: 0},
				{AgentType: "analyser", Position: 0},
			},
			wantErr: ErrDuplicatePosition,
		},
		{
			name:     "duplicate agent type",
			teamName: "my team",
			agents: []CreateAgentInput{
				{AgentType: "fetcher", Position: 0},
				{AgentType: "fetcher", Position: 1},
			},
			wantErr: ErrDuplicateAgentType,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, _ := newMockDB(t)
			svc := NewTeamService(db)

			_, _, err := svc.CreateTeam(context.Background(), uuid.New(), tc.teamName, tc.agents)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestCreateTeam_InvalidAgentTypeWrapsErr(t *testing.T) {
	db, _ := newMockDB(t)
	svc := NewTeamService(db)

	agents := []CreateAgentInput{{AgentType: "bad", Position: 0}}
	_, _, err := svc.CreateTeam(context.Background(), uuid.New(), "team", agents)

	if !errors.Is(err, ErrInvalidAgentType) {
		t.Errorf("expected ErrInvalidAgentType, got %v", err)
	}
}

func TestCreateTeam_DBError(t *testing.T) {
	db, mock := newMockDB(t)
	svc := NewTeamService(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teams"`)).WillReturnError(errors.New("db error"))
	mock.ExpectRollback()

	_, _, err := svc.CreateTeam(context.Background(), uuid.New(), "my team", baseAgents)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestCreateTeam_Success(t *testing.T) {
	db, mock := newMockDB(t)
	svc := NewTeamService(db)

	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teams"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "team_agents"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	team, agents, err := svc.CreateTeam(context.Background(), userID, "my team", baseAgents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, team.UserID)
	}
	if team.Name != "my team" {
		t.Errorf("expected name %q, got %q", "my team", team.Name)
	}
	if len(agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(agents))
	}
	if agents[0].AgentType != "fetcher" {
		t.Errorf("expected agent type fetcher, got %s", agents[0].AgentType)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestCreateTeam_NilContextDefaultsToEmptyMap(t *testing.T) {
	db, mock := newMockDB(t)
	svc := NewTeamService(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teams"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "team_agents"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	agents := []CreateAgentInput{{AgentType: "fetcher", Position: 0, Context: nil}}
	_, teamAgents, err := svc.CreateTeam(context.Background(), uuid.New(), "team", agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if teamAgents[0].Context == nil {
		t.Error("expected non-nil context map, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}
