package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "scaff",
		Short:        "A Go-native scaffolding and code-generation tool",
		Long:         "Scaff helps developers create and evolve Go projects with consistent structure and conventions.",
		SilenceUsage: true,
	}

	cmd.AddCommand(NewProjectCmd())
	return cmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
