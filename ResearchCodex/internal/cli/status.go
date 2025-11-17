package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jojo/ResearchCodex/internal/config"
	"github.com/jojo/ResearchCodex/internal/workspace"
	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the current rcodex status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, err := workspace.Detect()
			if err != nil {
				return err
			}
			cfg, err := config.Load(ws.ConfigPath())
			if err != nil {
				if errors.Is(err, config.ErrNotInitialized) {
					return errors.New("workspace not initialized. run rcodex init first")
				}
				return err
			}

			project := cfg.GetCurrentProject()
			if project == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "No project setted. please response by 'Please setup project by running rcodex project create <project_name>'")
				return nil
			}

			ideaPath := cfg.GetCurrentIdea()
			ideaSlug := filepath.Base(ideaPath)
			if ideaPath == "" || ideaSlug == "." {
				ideaSlug = "(none)"
			}
			mode := strings.ToLower(cfg.GetMode())
			if mode == "" {
				mode = "plan"
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "current_project: %s\n", project)
			fmt.Fprintf(out, "current_idea: %s\n", ideaSlug)
			fmt.Fprintf(out, "Mode: %s\n\n", strings.ToUpper(mode))

			switch mode {
			case "plan":
				return renderPlanStatus(out, ws, project, ideaSlug)
			case "code":
				return renderCodeStatus(out, ws, project, ideaSlug)
			default:
				fmt.Fprintln(out, "Mode is not set. Use rcodex plan or rcodex code to select a focus.")
				return nil
			}
		},
	}
}

func renderPlanStatus(out io.Writer, ws *workspace.Workspace, project, ideaSlug string) error {
	body, err := os.ReadFile(ws.PlanAgentsPath())
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "You are currently in PLAN mode")
	fmt.Fprintf(out, "- For PLAN mode, you will work in `projects/%s/%s`\n", project, ideaSlug)
	fmt.Fprintln(out, strings.TrimSpace(string(body)))
	return nil
}

func renderCodeStatus(out io.Writer, ws *workspace.Workspace, project, ideaSlug string) error {
	body, err := os.ReadFile(ws.CodeAgentsPath())
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "You are currently in CODE mode")
	fmt.Fprintf(out, "you will reference `projects/%s/%s` work on \n", project, ideaSlug)
	fmt.Fprintln(out, "  1. Shared code folder `srcs/`")
	fmt.Fprintf(out, "  2. Idea-specific scripts go in `experiments/%s/`.\n", project)
	fmt.Fprintln(out, strings.TrimSpace(string(body)))
	return nil
}
