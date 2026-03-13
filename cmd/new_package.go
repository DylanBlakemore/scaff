package cmd

import (
	"fmt"

	"scaff/internal/generator"
	pkggen "scaff/internal/project/pkg"

	"github.com/spf13/cobra"
)

type newPackageOpts struct {
	modulePath string
	style      string
	agents     bool
	makefile   bool
	ci         bool
	license    string
}

func NewPackageCmd() *cobra.Command {
	opts := &newPackageOpts{}

	cmd := &cobra.Command{
		Use:   "package <name>",
		Short: "Create a new Go package project",
		Long:  "Create a new standalone Go package project with optional tooling integrations.",
		Example: `  scaff new package mylib
  scaff new package mylib --module github.com/user/mylib
  scaff new package mylib --no-makefile --no-ci`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runNewPackage(args[0], opts)
		},
	}

	cmd.Flags().StringVar(&opts.modulePath, "module", "", "Go module path (defaults to project name)")
	cmd.Flags().StringVar(&opts.style, "style", "default", "Architecture style for the project")
	cmd.Flags().BoolVar(&opts.agents, "agents", true, "Generate an AGENTS.md file")
	cmd.Flags().BoolVar(&opts.makefile, "makefile", true, "Generate a Makefile")
	cmd.Flags().BoolVar(&opts.ci, "ci", true, "Generate CI configuration (GitHub Actions)")
	cmd.Flags().StringVar(&opts.license, "license", "none", "License SPDX ID to generate (e.g., MIT, Apache-2.0, BSD-3-Clause, MPL-2.0) or 'none'")

	return cmd
}

func runNewPackage(name string, opts *newPackageOpts) error {
	modulePath := opts.modulePath
	if modulePath == "" {
		modulePath = name
	}

	var features []string
	if opts.agents {
		features = append(features, "agents")
	}
	if opts.makefile {
		features = append(features, "makefile")
	}
	if opts.ci {
		features = append(features, "ci")
	}

	req := generator.Request{
		Name:       name,
		ModulePath: modulePath,
		Style:      opts.style,
		Features:   features,
		License:    opts.license,
	}

	gen := pkggen.New()
	result, err := gen.Generate(req)
	if err != nil {
		return err
	}

	fmt.Printf("Created package project %q in ./%s\n", modulePath, name)
	for _, f := range result.Files {
		fmt.Printf("  %s\n", f)
	}
	return nil
}
