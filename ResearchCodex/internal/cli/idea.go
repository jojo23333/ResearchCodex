package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func newIdeaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "idea",
		Short: "Work with ideas inside the current project",
	}
	cmd.AddCommand(newIdeaCreateCommand(), newIdeaStatusCommand())
	return cmd
}

func newIdeaCreateCommand() *cobra.Command {
	var body string
	var dependsOn string
	cmd := &cobra.Command{
		Use:   "create <title>",
		Short: "Create a new idea under the current project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			ws, err := workspace.Detect()
			if err != nil {
				return err
			}
			cfg, err := config.Load(ws.ConfigPath())
			if err != nil {
				return err
			}
			project := cfg.GetCurrentProject()
			if project == "" {
				return errors.New("no active project. run rcodex project create or rcodex project switch first")
			}

			now := time.Now().UTC()
			slug := fmt.Sprintf("%s_%s", timeutil.TimestampSlug(now), textutil.Slugify(title))
			relIdeaPath := filepath.ToSlash(filepath.Join("projects", project, slug))

			depPtr, err := determineDependency(ws, project, slug, dependsOn, cmd.Flags().Changed("depends-on"))
			if err != nil {
				return err
			}

			ideaDir := ws.ProjectIdeaDir(project, slug)
			if err := filesystem.EnsureDir(ideaDir); err != nil {
				return err
			}

			createdAt := timeutil.ISO8601(now)
			content := templates.IdeaMarkdown(title, createdAt, body)
			if err := filesystem.WriteFile(filepath.Join(ideaDir, "idea.md"), []byte(content)); err != nil {
				return err
			}
			if err := filesystem.WriteFile(filepath.Join(ideaDir, "plans.md"), []byte{}); err != nil {
				return err
			}

			entry := ideas.DependencyEntry{
				Project:   project,
				IdeaPath:  relIdeaPath,
				DependsOn: depPtr,
				CreatedAt: createdAt,
			}
			if err := ideas.AppendDependency(ws.IdeaDepsPath(), entry); err != nil {
				return err
			}

			cfg.SetCurrentIdea(relIdeaPath)
			cfg.SetMode("plan")
			if err := config.Save(ws.ConfigPath(), cfg); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created idea %s\n", relIdeaPath)
			return nil
		},
	}
	cmd.Flags().StringVar(&body, "body", "", "Idea body text")
	cmd.Flags().StringVar(&dependsOn, "depends-on", "", "Slug or path for dependency (use 'none' for a base idea)")
	return cmd
}

func determineDependency(ws *workspace.Workspace, project, newSlug, flagValue string, flagProvided bool) (*string, error) {
	projectDir := ws.ProjectDir(project)
	if flagProvided {
		if flagValue == "" || flagValue == "none" {
			return nil, nil
		}
		path, err := resolveIdeaPath(ws, project, flagValue)
		if err != nil {
			return nil, err
		}
		return &path, nil
	}

	latest, found, err := ideas.LatestIdeaDir(projectDir)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	if latest == newSlug {
		return nil, nil
	}
	path := filepath.ToSlash(filepath.Join("projects", project, latest))
	return &path, nil
}

func resolveIdeaPath(ws *workspace.Workspace, project, input string) (string, error) {
	if input == "" {
		return "", errors.New("dependency value cannot be empty")
	}

	if !strings.ContainsAny(input, "/\\") {
		slugPath := ws.ProjectIdeaDir(project, input)
		if stat, err := os.Stat(slugPath); err == nil && stat.IsDir() {
			return filepath.ToSlash(filepath.Join("projects", project, input)), nil
		}
	}

	var absPath string
	if filepath.IsAbs(input) {
		absPath = input
	} else {
		absPath = filepath.Join(ws.Root, filepath.FromSlash(input))
	}
	stat, err := os.Stat(absPath)
	if err != nil || !stat.IsDir() {
		return "", fmt.Errorf("could not find dependency %q", input)
	}
	rel, err := filepath.Rel(ws.Root, absPath)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("dependency path %q is outside the workspace", input)
	}
	return filepath.ToSlash(rel), nil
}

func newIdeaStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the current idea details and dependency chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, err := workspace.Detect()
			if err != nil {
				return err
			}
			cfg, err := config.Load(ws.ConfigPath())
			if err != nil {
				return err
			}
			ideaPath := cfg.GetCurrentIdea()
			if ideaPath == "" {
				return errors.New("no current idea selected. create or select an idea first")
			}

			absIdea := ws.Abs(ideaPath)
			ideaFile := filepath.Join(absIdea, "idea.md")
			title, created, body, err := parseIdeaFile(ideaFile)
			if err != nil {
				return err
			}
			project := cfg.GetCurrentProject()
			if project == "" {
				project = deriveProjectFromIdeaPath(ideaPath)
			}

			deps, err := ideas.LoadDependencies(ws.IdeaDepsPath())
			if err != nil {
				return err
			}
			chain := ideas.ResolveChain(deps, ideaPath)

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Idea: %s\n", titleOrFallback(title, ideaPath))
			fmt.Fprintf(out, "Project: %s\n", projectOrFallback(project))
			fmt.Fprintf(out, "Created: %s\n", created)
			fmt.Fprintf(out, "Path: %s\n", ideaPath)
			fmt.Fprintf(out, "Dependency chain:\n")
			if len(chain) == 0 {
				fmt.Fprintln(out, "  (no dependency records)")
			} else {
				for i, entry := range chain {
					prefix := "  "
					if i > 0 {
						prefix = "  -> "
					}
					fmt.Fprintf(out, "%s%s\n", prefix, entry.IdeaPath)
				}
			}
			fmt.Fprintf(out, "\nBody:\n%s\n", indentBody(body))
			return nil
		},
	}
}

func parseIdeaFile(path string) (title string, created string, body string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inBody := false
	var bodyLines []string
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "# ") && title == "":
			title = strings.TrimSpace(line[2:])
		case strings.HasPrefix(line, "Created:") && created == "":
			created = strings.TrimSpace(strings.TrimPrefix(line, "Created:"))
		case strings.HasPrefix(line, "Body:"):
			inBody = true
		default:
			if inBody {
				bodyLines = append(bodyLines, line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", "", err
	}
	return title, created, strings.TrimSpace(strings.Join(bodyLines, "\n")), nil
}

func deriveProjectFromIdeaPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

func titleOrFallback(title, ideaPath string) string {
	if title != "" {
		return title
	}
	return filepath.Base(ideaPath)
}

func projectOrFallback(project string) string {
	if project != "" {
		return project
	}
	return "(unknown)"
}

func indentBody(body string) string {
	if strings.TrimSpace(body) == "" {
		return "  (empty)"
	}
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		lines[i] = "  " + line
	}
	return strings.Join(lines, "\n")
}
