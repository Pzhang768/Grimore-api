package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Pzhang768/Grimore-api/internal/agents"
)

func TestNewAgentHandler(t *testing.T) {
	h := NewAgentHandler()
	if h == nil {
		t.Error("expected non-nil handler")
	}
}

func TestAgentRegister(t *testing.T) {
	r := gin.New()
	NewAgentHandler().Register(r)
	routes := r.Routes()
	if len(routes) != 1 || routes[0].Method != http.MethodGet || routes[0].Path != "/agents" {
		t.Errorf("unexpected routes: %+v", routes)
	}
}

func TestListAgents_ReturnsAllAgents(t *testing.T) {
	r := gin.New()
	NewAgentHandler().Register(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/agents", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []agentDefinition
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp) != len(agents.Registry) {
		t.Errorf("expected %d agents, got %d", len(agents.Registry), len(resp))
	}

	for _, def := range resp {
		if def.AgentType == "" {
			t.Error("agent definition has empty agent_type")
		}
		if def.Description == "" {
			t.Error("agent definition has empty description")
		}
		if def.ContextSchema == nil {
			t.Errorf("agent %q has nil context_schema", def.AgentType)
		}
	}
}
