package cmd

import (
	nebula "github.com/nihcosnosaj/nebula-homelab"
	"github.com/nihcosnosaj/nebula-homelab/internal/cluster"
	"github.com/nihcosnosaj/nebula-homelab/internal/platform"
	"github.com/spf13/cobra"
)

var dryRun bool

func newUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Spin up the Nebula k3s cluster",
		Long:  `Nebula 'up' provisions AWS spot instances via Terraform and configures k3s via Ansible.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// initializae providers
			tf := platform.NewTerraformExec("./terraform", nebula.ProjectFiles)
			ans := platform.NewAnsibleExec("./ansible/inventory.ini", nebula.ProjectFiles)

			manager := cluster.Manager{
				TF:      tf,
				Ansible: ans,
			}

			return manager.Up(dryRun)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview infrastructure changes without applying them.")
	return cmd
}
