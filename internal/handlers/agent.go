package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Pzhang768/Grimore-api/internal/agents"
)

type agentDefinition struct {
	AgentType     string                `json:"agent_type"`
	Description   string                `json:"description"`
	ContextSchema []agents.ContextField `json:"context_schema"`
}

type AgentHandler struct{}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{}
}

func (h *AgentHandler) Register(r gin.IRouter) {
	r.GET("/agents", h.ListAgents)
}

func (h *AgentHandler) ListAgents(c *gin.Context) {
	resp := make([]agentDefinition, 0, len(agents.Registry))
	for _, a := range agents.Registry {
		resp = append(resp, agentDefinition{
			AgentType:     a.Type(),
			Description:   a.Description(),
			ContextSchema: a.ContextSchema(),
		})
	}
	c.JSON(http.StatusOK, resp)
}
