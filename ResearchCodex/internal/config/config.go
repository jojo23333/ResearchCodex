package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ErrNotInitialized indicates the workspace has not been bootstrapped with rcodex init.
var ErrNotInitialized = errors.New("rcodex not initialized. Please run rcodex init first")

// Config mirrors .rcodex/config.yaml.
type Config struct {
	CurrentProject *string `yaml:"current_project"`
	CurrentIdea    *string `yaml:"current_idea"`
	Mode           *string `yaml:"mode"`
}

func Default() *Config {
	return &Config{}
}

func (c *Config) GetCurrentProject() string {
	if c.CurrentProject == nil {
		return ""
	}
	return *c.CurrentProject
}

func (c *Config) GetCurrentIdea() string {
	if c.CurrentIdea == nil {
		return ""
	}
	return *c.CurrentIdea
}

func (c *Config) GetMode() string {
	if c.Mode == nil {
		return ""
	}
	return *c.Mode
}

func (c *Config) SetCurrentProject(project string) {
	c.CurrentProject = strPtr(project)
}

func (c *Config) ClearCurrentIdea() {
	c.CurrentIdea = nil
}

func (c *Config) SetCurrentIdea(relPath string) {
	c.CurrentIdea = strPtr(relPath)
}

func (c *Config) SetMode(mode string) {
	c.Mode = strPtr(mode)
}

func (c *Config) ClearMode() {
	c.Mode = nil
}

func strPtr(v string) *string {
	return &v
}

// Load reads config.yaml from disk.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotInitialized
		}
		return nil, err
	}
	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save persists config.yaml, creating directories as needed.
func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
