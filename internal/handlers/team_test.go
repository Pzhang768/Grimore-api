package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Pzhang768/Grimore-api/internal/middleware"
	"github.com/Pzhang768/Grimore-api/internal/models"
	"github.com/Pzhang768/Grimore-api/internal/services"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type fakeTeamCreator struct {
	team   *models.Team
	agents []models.TeamAgent
	err    error
}

func (f *fakeTeamCreator) CreateTeam(_ context.Context, _ uuid.UUID, _ string, _ []services.CreateAgentInput) (*models.Team, []models.TeamAgent, error) {
	return f.team, f.agents, f.err
}

func newTeamRouter(fake teamCreator, userID uuid.UUID) *gin.Engine {
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, userID)
		c.Next()
	})
	h := &TeamHandler{svc: fake}
	r.POST("/teams", h.CreateTeam)
	return r
}

func postTeams(r *gin.Engine, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func TestNewTeamHandler(t *testing.T) {
	fake := &fakeTeamCreator{}
	h := NewTeamHandler(nil)
	h.svc = fake
	if h.svc == nil {
		t.Error("expected non-nil svc")
	}
}

func TestRegister(t *testing.T) {
	h := &TeamHandler{svc: &fakeTeamCreator{}}
	r := gin.New()
	h.Register(r)
	routes := r.Routes()
	if len(routes) != 1 || routes[0].Method != http.MethodPost || routes[0].Path != "/teams" {
		t.Errorf("unexpected routes: %+v", routes)
	}
}

func TestCreateTeamHandler_InvalidJSON(t *testing.T) {
	r := newTeamRouter(&fakeTeamCreator{}, uuid.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/teams", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTeamHandler_MissingRequiredFields(t *testing.T) {
	r := newTeamRouter(&fakeTeamCreator{}, uuid.New())

	w := postTeams(r, map[string]interface{}{
		"agents": []map[string]interface{}{{"agent_type": "fetcher", "position": 0}},
	})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTeamHandler_ValidationErrors_ReturnBadRequest(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"name required", services.ErrTeamNameRequired},
		{"agents required", services.ErrAgentsRequired},
		{"too many agents", services.ErrTooManyAgents},
		{"invalid agent type", services.ErrInvalidAgentType},
		{"duplicate position", services.ErrDuplicatePosition},
		{"duplicate agent type", services.ErrDuplicateAgentType},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := newTeamRouter(&fakeTeamCreator{err: tc.err}, uuid.New())

			w := postTeams(r, map[string]interface{}{
				"name":   "team",
				"agents": []map[string]interface{}{{"agent_type": "fetcher", "position": 0}},
			})

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", w.Code)
			}
		})
	}
}

func TestCreateTeamHandler_InternalError(t *testing.T) {
	r := newTeamRouter(&fakeTeamCreator{err: errors.New("db exploded")}, uuid.New())

	w := postTeams(r, map[string]interface{}{
		"name":   "team",
		"agents": []map[string]interface{}{{"agent_type": "fetcher", "position": 0}},
	})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestCreateTeamHandler_Success(t *testing.T) {
	teamID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	fake := &fakeTeamCreator{
		team: &models.Team{
			ID:        teamID,
			UserID:    userID,
			Name:      "dream team",
			CreatedAt: now,
		},
		agents: []models.TeamAgent{
			{TeamID: teamID, AgentType: "fetcher", Position: 0, Context: map[string]any{}},
		},
	}
	r := newTeamRouter(fake, userID)

	w := postTeams(r, map[string]interface{}{
		"name":   "dream team",
		"agents": []map[string]interface{}{{"agent_type": "fetcher", "position": 0}},
	})

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var resp createTeamResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != teamID {
		t.Errorf("expected team ID %s, got %s", teamID, resp.ID)
	}
	if resp.Name != "dream team" {
		t.Errorf("expected name %q, got %q", "dream team", resp.Name)
	}
	if len(resp.Agents) != 1 || resp.Agents[0].AgentType != "fetcher" {
		t.Errorf("unexpected agents: %+v", resp.Agents)
	}
}
