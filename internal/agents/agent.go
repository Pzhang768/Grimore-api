package agents

import "context"

type ContextField struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
}

type AgentInput struct {
	Context map[string]any
	Payload string
}

type AgentOutput struct {
	Content string
}

type Agent interface {
	Type() string
	Description() string
	ContextSchema() []ContextField
	Run(ctx context.Context, input AgentInput) (AgentOutput, error)
}
