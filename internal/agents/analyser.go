package agents

import "context"

type Analyser struct{}

func (a *Analyser) Type() string        { return "analyser" }
func (a *Analyser) Description() string { return "Scores fit between job listings and your resume" }
func (a *Analyser) ContextSchema() []ContextField {
	return []ContextField{
		{Key: "min_fit_score", Label: "Minimum fit score (0–100)", Required: false},
	}
}
func (a *Analyser) Run(_ context.Context, _ AgentInput) (AgentOutput, error) {
	return AgentOutput{Content: "stub: analyser not yet implemented"}, nil
}
