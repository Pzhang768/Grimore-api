package agents

var Registry = map[string]Agent{
	"fetcher":     &Fetcher{},
	"analyser":    &Analyser{},
	"tailor":      &Tailor{},
	"coordinator": &Coordinator{},
}
