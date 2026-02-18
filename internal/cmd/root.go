package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "nebula",
		Short: "Nebula: Jason's All-Around Great Homelab Manager :) ",
		Long:  `A robust CLI for managing k3s clusters on AWS Spot Instances using Terraform for provisioning and Ansible for configuration.`,
	}

	// available to all sub-commands: verbose, todo: dryrun
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "show detailed output")

	// subcommand registry
	rootCmd.AddCommand(newUpCmd())
	rootCmd.AddCommand(newDestroyCmd())
	rootCmd.AddCommand(newStatusCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
