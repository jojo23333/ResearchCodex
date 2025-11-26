package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jojo/ResearchCodex/internal/config"
	"github.com/jojo/ResearchCodex/internal/templates"
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
				mode = "scope"
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "current_project: %s\n", project)
			fmt.Fprintf(out, "current_scope: %s\n", ideaSlug)
			fmt.Fprintf(out, "Mode: %s\n\n", strings.ToUpper(mode))

			switch mode {
			case "scope":
				return renderScopeStatus(out, ws, project, ideaSlug)
			case "plan":
				return renderPlanStatus(out, ws, project, ideaSlug)
			case "code":
				return renderCodeStatus(out, ws, project, ideaSlug)
			default:
				fmt.Fprintln(out, "Mode is not set. Use rcodex scope, rcodex plan or rcodex code to select a focus.")
				return nil
			}
		},
	}
}

func renderScopeStatus(out io.Writer, ws *workspace.Workspace, project, ideaSlug string) error {
	body, err := readAgentFile(ws.ScopeAgentsPath(), templates.ScopeModeAgentsMarkdown())
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "You are currently in SCOPE mode")
	fmt.Fprintf(out, "- For SCOPE mode, you will work in `projects/%s/%s`\n", project, ideaSlug)
	fmt.Fprintln(out, strings.TrimSpace(string(body)))
	return nil
}

func renderPlanStatus(out io.Writer, ws *workspace.Workspace, project, ideaSlug string) error {
	body, err := readAgentFile(ws.PlanAgentsPath(), templates.PlanModeAgentsMarkdown())
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "You are currently in PLAN mode")
	fmt.Fprintf(out, "- For PLAN mode, you will work in `projects/%s/%s`\n", project, ideaSlug)
	fmt.Fprintln(out, strings.TrimSpace(string(body)))
	return nil
}

func renderCodeStatus(out io.Writer, ws *workspace.Workspace, project, ideaSlug string) error {
	body, err := readAgentFile(ws.CodeAgentsPath(), templates.CodeModeAgentsMarkdown())
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

func readAgentFile(path, fallback string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fallback, nil
		}
		return "", err
	}
	return string(data), nil
}
