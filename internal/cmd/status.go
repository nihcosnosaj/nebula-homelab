package cmd

import (
	"github.com/nihcosnosaj/nebula-homelab/internal/cluster"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "View the current state of Nebula instances in AWS",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := &cluster.Manager{}
			return manager.Status("nebula")
		},
	}
}
