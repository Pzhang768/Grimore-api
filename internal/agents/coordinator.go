package agents

import "context"

type Coordinator struct{}

func (c *Coordinator) Type() string        { return "coordinator" }
func (c *Coordinator) Description() string { return "Passes context between agents and signals completion" }
func (c *Coordinator) ContextSchema() []ContextField {
	return []ContextField{
		{Key: "max_iterations", Label: "Max iterations (default 3)", Required: false},
	}
}
func (c *Coordinator) Run(_ context.Context, _ AgentInput) (AgentOutput, error) {
	return AgentOutput{Content: "stub: coordinator not yet implemented"}, nil
}
