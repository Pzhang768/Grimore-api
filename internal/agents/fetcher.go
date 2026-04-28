package agents

import "context"

type Fetcher struct{}

func (f *Fetcher) Type() string        { return "fetcher" }
func (f *Fetcher) Description() string { return "Finds relevant job listings based on role and location" }
func (f *Fetcher) ContextSchema() []ContextField {
	return []ContextField{
		{Key: "role_title", Label: "Role title", Required: true},
		{Key: "location", Label: "Location", Required: false},
		{Key: "keywords", Label: "Keywords", Required: false},
	}
}
func (f *Fetcher) Run(_ context.Context, _ AgentInput) (AgentOutput, error) {
	return AgentOutput{Content: "stub: fetcher not yet implemented"}, nil
}
