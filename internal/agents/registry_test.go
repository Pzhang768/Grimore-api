package agents

import (
	"context"
	"testing"
)

func TestRegistry_ContainsAllAgentTypes(t *testing.T) {
	expected := []string{"fetcher", "analyser", "tailor", "coordinator"}
	for _, agentType := range expected {
		if _, ok := Registry[agentType]; !ok {
			t.Errorf("registry missing agent type %q", agentType)
		}
	}
}

func TestRegistry_TypeMatchesKey(t *testing.T) {
	for key, a := range Registry {
		if a.Type() != key {
			t.Errorf("agent key %q has Type() %q — they must match", key, a.Type())
		}
	}
}

func TestRegistry_AllAgentsHaveDescription(t *testing.T) {
	for _, a := range Registry {
		if a.Description() == "" {
			t.Errorf("agent %q has empty description", a.Type())
		}
	}
}

func TestRegistry_AllAgentsHaveContextSchema(t *testing.T) {
	for _, a := range Registry {
		schema := a.ContextSchema()
		if schema == nil {
			t.Errorf("agent %q returned nil context schema", a.Type())
		}
		for _, field := range schema {
			if field.Key == "" || field.Label == "" {
				t.Errorf("agent %q has context field with empty key or label: %+v", a.Type(), field)
			}
		}
	}
}

func TestRegistry_AllAgentsRun(t *testing.T) {
	for _, a := range Registry {
		out, err := a.Run(context.Background(), AgentInput{})
		if err != nil {
			t.Errorf("agent %q Run() returned error: %v", a.Type(), err)
		}
		if out.Content == "" {
			t.Errorf("agent %q Run() returned empty content", a.Type())
		}
	}
}
