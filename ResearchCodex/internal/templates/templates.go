package templates

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed CodeModeAgents.md PlanModeAgents.md RootAgents.md
var templateFS embed.FS

var (
	planModeAgents = mustReadTemplate("PlanModeAgents.md")
	codeModeAgents = mustReadTemplate("CodeModeAgents.md")
	rootAgents     = mustReadTemplate("RootAgents.md")
)

func mustReadTemplate(path string) string {
	data, err := templateFS.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("templates: missing %s: %w", path, err))
	}
	return string(data)
}

func PlanModeAgentsMarkdown() string { return planModeAgents }

func CodeModeAgentsMarkdown() string { return codeModeAgents }

func RootAgentsMarkdown() string { return rootAgents }

const DefaultProjectAgents = `# Project Agents

- Planner: refresh context by reading projects/<project>/ideas
- Executor: implements todos derived from the active idea
- Analyst: interprets experiment results and records outcomes
`

// IdeaMarkdown builds the canonical content for idea.md.
func IdeaMarkdown(title, createdAt, body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		body = "(fill this in with hypothesis, questions, and outline)"
	}

	return fmt.Sprintf(`# %s

Created: %s

Body:
%s
`, title, createdAt, body)
}

// PlansMarkdown returns the scaffold for plans.md.
func PlansMarkdown() string {
	return "# Plans\n\n"
}
