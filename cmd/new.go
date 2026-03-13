package cmd

import "github.com/spf13/cobra"

func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new project",
		Long:  "Create a new Go project from a supported project type.",
	}
	cmd.AddCommand(NewPackageCmd())
	return cmd
}
