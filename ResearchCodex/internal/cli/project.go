package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jojo/ResearchCodex/internal/config"
	"github.com/jojo/ResearchCodex/internal/filesystem"
	"github.com/jojo/ResearchCodex/internal/ideas"
	"github.com/jojo/ResearchCodex/internal/templates"
	"github.com/jojo/ResearchCodex/internal/textutil"
	"github.com/jojo/ResearchCodex/internal/timeutil"
	"github.com/jojo/ResearchCodex/internal/workspace"
	"github.com/spf13/cobra"
)

func newProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage ResearchCodex projects",
	}
	cmd.AddCommand(newProjectNewCommand(), newProjectSwitchCommand())
	return cmd
}

func newProjectNewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "new <name>",
		Short:   "Create a new research project",
		Aliases: []string{"create"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if name == "" {
				return errors.New("project name cannot be empty")
			}

			ws, err := workspace.Detect()
			if err != nil {
				return err
			}
			if _, err := os.Stat(ws.RCodexDir()); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return errors.New("workspace not initialized. run rcodex init first")
				}
				return err
			}

			projectDir := ws.ProjectDir(name)
			if _, err := os.Stat(projectDir); err == nil {
				return fmt.Errorf("project %q already exists", name)
			}

			if err := filesystem.EnsureDir(projectDir); err != nil {
				return err
			}
			if err := filesystem.EnsureDir(ws.ProjectContextDir(name)); err != nil {
				return err
			}
			if err := filesystem.WriteFile(ws.ProjectAgentsPath(name), []byte(templates.DefaultProjectAgents)); err != nil {
				return err
			}
			if err := filesystem.EnsureDir(ws.ProjectExperimentsDir(name)); err != nil {
				return err
			}

			slugPart, createdAt := createBaseIdeaSlug(name)
			ideaDir := ws.ProjectIdeaDir(name, slugPart)
			if err := filesystem.EnsureDir(ideaDir); err != nil {
				return err
			}

			ideaContent := templates.IdeaMarkdown(name, createdAt, "")
			if err := filesystem.WriteFile(filepath.Join(ideaDir, "idea.md"), []byte(ideaContent)); err != nil {
				return err
			}
			if err := filesystem.WriteFile(filepath.Join(ideaDir, "plans.md"), []byte{}); err != nil {
				return err
			}

			relIdeaPath := filepath.ToSlash(filepath.Join("projects", name, slugPart))
			entry := ideas.DependencyEntry{
				Project:   name,
				IdeaPath:  relIdeaPath,
				DependsOn: nil,
				CreatedAt: createdAt,
			}
			if err := ideas.AppendDependency(ws.IdeaDepsPath(), entry); err != nil {
				return err
			}

			cfg, err := config.Load(ws.ConfigPath())
			if err != nil {
				if errors.Is(err, config.ErrNotInitialized) {
					cfg = config.Default()
				} else {
					return err
				}
			}
			cfg.SetCurrentProject(name)
			cfg.SetCurrentIdea(relIdeaPath)
			cfg.SetMode("scope")
			if err := config.Save(ws.ConfigPath(), cfg); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created project %q with base scope %s\n", name, relIdeaPath)
			return nil
		},
	}
}

func newProjectSwitchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <name>",
		Short: "Switch the active project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			ws, err := workspace.Detect()
			if err != nil {
				return err
			}

			cfg, err := config.Load(ws.ConfigPath())
			if err != nil {
				return err
			}

			if _, err := os.Stat(ws.ProjectDir(name)); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("project %q does not exist", name)
				}
				return err
			}

			cfg.SetCurrentProject(name)
			cfg.ClearCurrentIdea()
			if err := config.Save(ws.ConfigPath(), cfg); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Switched to project %q\n", name)
			return nil
		},
	}
}

func createBaseIdeaSlug(projectName string) (slug string, createdAt string) {
	now := time.Now().UTC()
	return fmt.Sprintf("%s_%s", timeutil.TimestampSlug(now), textutil.Slugify(projectName)), timeutil.ISO8601(now)
}
