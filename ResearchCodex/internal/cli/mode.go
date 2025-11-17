package cli

import (
	"errors"
	"fmt"

	"github.com/jojo/ResearchCodex/internal/config"
	"github.com/jojo/ResearchCodex/internal/workspace"
	"github.com/spf13/cobra"
)

func newPlanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Switch to PLAN mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setMode(cmd, "plan")
		},
	}
}

func newCodeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "code",
		Short: "Switch to CODE mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			return setMode(cmd, "code")
		},
	}
}

func setMode(cmd *cobra.Command, mode string) error {
	ws, err := workspace.Detect()
	if err != nil {
		return err
	}
	cfg, err := config.Load(ws.ConfigPath())
	if err != nil {
		return err
	}
	if cfg.GetCurrentProject() == "" {
		return errors.New("no project set. create or switch to a project first")
	}
	cfg.SetMode(mode)
	if err := config.Save(ws.ConfigPath(), cfg); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Switched to %s mode\n", mode)
	return nil
}
