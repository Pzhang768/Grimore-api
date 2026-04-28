package agents

import "context"

type Tailor struct{}

func (t *Tailor) Type() string        { return "tailor" }
func (t *Tailor) Description() string { return "Rewrites your resume to match each job listing" }
func (t *Tailor) ContextSchema() []ContextField {
	return []ContextField{
		{Key: "tone", Label: "Tone (e.g. professional, concise)", Required: false},
	}
}
func (t *Tailor) Run(_ context.Context, _ AgentInput) (AgentOutput, error) {
	return AgentOutput{Content: "stub: tailor not yet implemented"}, nil
}
