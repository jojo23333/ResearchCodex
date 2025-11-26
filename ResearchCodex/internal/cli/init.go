package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jojo/ResearchCodex/internal/config"
	"github.com/jojo/ResearchCodex/internal/filesystem"
	"github.com/jojo/ResearchCodex/internal/templates"
	"github.com/jojo/ResearchCodex/internal/workspace"
	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new ResearchCodex workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, err := workspace.Detect()
			if err != nil {
				return err
			}

			if _, err := os.Stat(ws.RCodexDir()); err == nil {
				return fmt.Errorf("rcodex already initialized (found %s)", ws.RCodexDir())
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}

			if err := filesystem.EnsureDir(ws.RCodexDir()); err != nil {
				return err
			}
			for _, dir := range []string{ws.ProjectsDir(), ws.ExperimentsDir(), ws.SrcsDir()} {
				if err := filesystem.EnsureDir(dir); err != nil {
					return err
				}
			}

			if err := filesystem.WriteFile(ws.PlanAgentsPath(), []byte(templates.PlanModeAgentsMarkdown())); err != nil {
				return err
			}
			if err := filesystem.WriteFile(ws.CodeAgentsPath(), []byte(templates.CodeModeAgentsMarkdown())); err != nil {
				return err
			}
			if err := filesystem.WriteFile(ws.ScopeAgentsPath(), []byte(templates.ScopeModeAgentsMarkdown())); err != nil {
				return err
			}
			if err := filesystem.WriteFile(ws.IdeaDepsPath(), []byte{}); err != nil {
				return err
			}
			if err := filesystem.WriteFile(ws.RootAgentsPath(), []byte(templates.RootAgentsMarkdown())); err != nil {
				return err
			}

			projectsAgents := filepath.Join(ws.ProjectsDir(), "AGENTS.md")
			if err := filesystem.WriteFile(projectsAgents, []byte(templates.DefaultProjectAgents)); err != nil {
				return err
			}

			if err := config.Save(ws.ConfigPath(), config.Default()); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Initialized ResearchCodex workspace.")
			return nil
		},
	}
}
