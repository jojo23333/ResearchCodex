package workspace

import (
	"os"
	"path/filepath"
)

// Workspace represents the current working tree root where rcodex commands run.
type Workspace struct {
	Root string
}

// Detect resolves the current working directory as the workspace root.
func Detect() (*Workspace, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Workspace{Root: root}, nil
}

func (w *Workspace) RCodexDir() string {
	return filepath.Join(w.Root, ".rcodex")
}

func (w *Workspace) ConfigPath() string {
	return filepath.Join(w.RCodexDir(), "config.yaml")
}

func (w *Workspace) IdeaDepsPath() string {
	return filepath.Join(w.RCodexDir(), "idea_deps.jsonl")
}

func (w *Workspace) PlanAgentsPath() string {
	return filepath.Join(w.RCodexDir(), ".plan_mode_agents.md")
}

func (w *Workspace) CodeAgentsPath() string {
	return filepath.Join(w.RCodexDir(), ".code_mode_agents.md")
}

func (w *Workspace) ScopeAgentsPath() string {
	return filepath.Join(w.RCodexDir(), ".scope_mode_agents.md")
}

func (w *Workspace) ProjectsDir() string {
	return filepath.Join(w.Root, "projects")
}

func (w *Workspace) RootAgentsPath() string {
	return filepath.Join(w.Root, "AGENTS.md")
}

func (w *Workspace) ExperimentsDir() string {
	return filepath.Join(w.Root, "experiments")
}

func (w *Workspace) SrcsDir() string {
	return filepath.Join(w.Root, "srcs")
}

func (w *Workspace) ProjectDir(name string) string {
	return filepath.Join(w.ProjectsDir(), name)
}

func (w *Workspace) ProjectContextDir(name string) string {
	return filepath.Join(w.ProjectDir(name), "context")
}

func (w *Workspace) ProjectAgentsPath(name string) string {
	return filepath.Join(w.ProjectDir(name), "AGENTS.md")
}

func (w *Workspace) ProjectExperimentsDir(name string) string {
	return filepath.Join(w.ExperimentsDir(), name)
}

func (w *Workspace) ProjectIdeaDir(project, slug string) string {
	return filepath.Join(w.ProjectDir(project), slug)
}

// Abs resolves a workspace-relative path (using '/' separators) to an OS path.
func (w *Workspace) Abs(rel string) string {
	if rel == "" {
		return w.Root
	}
	return filepath.Join(w.Root, filepath.FromSlash(rel))
}
