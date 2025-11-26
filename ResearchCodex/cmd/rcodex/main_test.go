package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type configFile struct {
	CurrentProject *string `yaml:"current_project"`
	CurrentIdea    *string `yaml:"current_idea"`
	Mode           *string `yaml:"mode"`
}

func TestEndToEndWorkflow(t *testing.T) {
	root := repoRoot(t)
	tmp := t.TempDir()
	bin := filepath.Join(tmp, "rcodex")

	goCache := filepath.Join(tmp, "gocache")
	if err := os.MkdirAll(goCache, 0o755); err != nil {
		t.Fatalf("mkdir gocache: %v", err)
	}

	build := exec.Command("go", "build", "-o", bin, "./cmd/rcodex")
	build.Dir = root
	build.Env = append(os.Environ(), "GOCACHE="+goCache)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build rcodex: %v\n%s", err, out)
	}

	run := func(name string, args ...string) string {
		cmd := exec.Command(bin, args...)
		cmd.Dir = tmp
		cmd.Env = append(os.Environ(), "GOCACHE="+goCache)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s failed: %v\n%s", strings.Join(args, " "), err, out)
		}
		return string(out)
	}

	run("init", "init")
	assertPathExists(t, filepath.Join(tmp, ".rcodex"))

	// project new should set scope mode and seed base scope using project name
	run("project new", "project", "new", "alpha")

	cfg := readConfig(t, filepath.Join(tmp, ".rcodex", "config.yaml"))
	if got, want := deref(cfg.CurrentProject), "alpha"; got != want {
		t.Fatalf("current_project: got %q want %q", got, want)
	}
	if got, want := deref(cfg.Mode), "scope"; got != want {
		t.Fatalf("mode: got %q want %q", got, want)
	}
	baseScope := deref(cfg.CurrentIdea)
	if !strings.Contains(baseScope, "alpha") {
		t.Fatalf("current_scope path %q does not contain project name", baseScope)
	}
	if !regexp.MustCompile(`projects/alpha/\d{8}_\d{6}_alpha`).MatchString(baseScope) {
		t.Fatalf("current_scope path %q does not match expected slug pattern", baseScope)
	}

	// status should mention scope mode and current project
	statusOut := run("status", "status")
	if !strings.Contains(statusOut, "SCOPE mode") {
		t.Fatalf("status missing SCOPE mode: %s", statusOut)
	}
	if !strings.Contains(statusOut, "current_project: alpha") {
		t.Fatalf("status missing project: %s", statusOut)
	}

	// create a second scope; should depend on the first
	run("idea new", "idea", "new", "Second scope", "--body", "demo body")

	cfg = readConfig(t, filepath.Join(tmp, ".rcodex", "config.yaml"))
	secondScope := deref(cfg.CurrentIdea)
	if secondScope == baseScope {
		t.Fatalf("second scope path matches base scope")
	}

	deps := readFile(t, filepath.Join(tmp, ".rcodex", "idea_deps.jsonl"))
	if !strings.Contains(deps, `"depends_on":"`+baseScope+`"`) {
		t.Fatalf("dependency log missing link to base scope:\n%s", deps)
	}

	// mode transitions
	run("plan", "plan")
	cfg = readConfig(t, filepath.Join(tmp, ".rcodex", "config.yaml"))
	if got := deref(cfg.Mode); got != "plan" {
		t.Fatalf("mode after plan: got %q", got)
	}

	run("code", "code")
	cfg = readConfig(t, filepath.Join(tmp, ".rcodex", "config.yaml"))
	if got := deref(cfg.Mode); got != "code" {
		t.Fatalf("mode after code: got %q", got)
	}

	ideaStatus := run("idea status", "idea", "status")
	if !strings.Contains(ideaStatus, "Scope:") {
		t.Fatalf("idea status missing Scope header: %s", ideaStatus)
	}
	if !strings.Contains(ideaStatus, baseScope) {
		t.Fatalf("idea status missing dependency chain entry: %s", ideaStatus)
	}
}

// repoRoot walks up from the current working directory to find go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}

func assertPathExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected path %s to exist: %v", path, err)
	}
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func readConfig(t *testing.T, path string) configFile {
	t.Helper()
	data := readFile(t, path)
	var cfg configFile
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}
	return cfg
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}
	return string(bytes.TrimSpace(b))
}
