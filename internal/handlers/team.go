package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Pzhang768/Grimore-api/internal/middleware"
	"github.com/Pzhang768/Grimore-api/internal/models"
	"github.com/Pzhang768/Grimore-api/internal/services"
)

type agentInput struct {
	AgentType string         `json:"agent_type" binding:"required"`
	Position  int            `json:"position"`
	Context   map[string]any `json:"context"`
}

type createTeamRequest struct {
	Name   string       `json:"name" binding:"required"`
	Agents []agentInput `json:"agents" binding:"required"`
}

type agentResponse struct {
	AgentType string         `json:"agent_type"`
	Position  int            `json:"position"`
	Context   map[string]any `json:"context"`
}

type createTeamResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Agents    []agentResponse `json:"agents"`
	CreatedAt time.Time       `json:"created_at"`
}

type teamCreator interface {
	CreateTeam(ctx context.Context, userID uuid.UUID, name string, agents []services.CreateAgentInput) (*models.Team, []models.TeamAgent, error)
}

type TeamHandler struct {
	svc teamCreator
}

func NewTeamHandler(svc *services.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) Register(r gin.IRouter) {
	r.POST("/teams", h.CreateTeam)
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req createTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic}) //nolint:errcheck
		return
	}

	// Auth middleware always sets ContextKeyUserID before this handler is reached.
	userID, _ := c.Get(middleware.ContextKeyUserID)

	inputs := make([]services.CreateAgentInput, len(req.Agents))
	for i, a := range req.Agents {
		inputs[i] = services.CreateAgentInput{
			AgentType: a.AgentType,
			Position:  a.Position,
			Context:   a.Context,
		}
	}

	team, agents, err := h.svc.CreateTeam(c.Request.Context(), userID.(uuid.UUID), req.Name, inputs)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTeamNameRequired),
			errors.Is(err, services.ErrAgentsRequired),
			errors.Is(err, services.ErrTooManyAgents),
			errors.Is(err, services.ErrInvalidAgentType),
			errors.Is(err, services.ErrDuplicatePosition),
			errors.Is(err, services.ErrDuplicateAgentType):
			c.Error(&gin.Error{Err: err, Type: gin.ErrorTypePublic}) //nolint:errcheck
		default:
			c.Error(err) //nolint:errcheck
		}
		return
	}

	resp := createTeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		CreatedAt: team.CreatedAt,
		Agents:    make([]agentResponse, len(agents)),
	}
	for i, a := range agents {
		resp.Agents[i] = agentResponse{
			AgentType: a.AgentType,
			Position:  a.Position,
			Context:   a.Context,
		}
	}

	c.JSON(http.StatusCreated, resp)
}
