package internal

import "github.com/spf13/cobra"

func init() {
	GenerateCmd.AddCommand(
		ModuleCmd,
	)
}

// GenerateCmd generates boilerplate code with different commands
var GenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generate boilerplate code with different commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
