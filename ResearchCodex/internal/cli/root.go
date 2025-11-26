package cli

import "github.com/spf13/cobra"

// NewRootCommand wires up the rcodex CLI root entry point.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "rcodex",
		Short:         "ResearchCodex CLI",
		Long:          "ResearchCodex: lightweight research planning CLI for coding agents.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		newInitCommand(),
		newProjectCommand(),
		newIdeaCommand(),
		newStatusCommand(),
		newScopeModeCommand(),
		newPlanCommand(),
		newCodeCommand(),
	)

	return root
}
