package cmd

import (
	"fmt"
	"strings"

	nebula "github.com/nihcosnosaj/nebula-homelab"
	"github.com/nihcosnosaj/nebula-homelab/internal/cluster"
	"github.com/nihcosnosaj/nebula-homelab/internal/platform"
	"github.com/spf13/cobra"
)

var forceDestroy bool

func newDestroyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Tear down the Nebula cluster and AWS resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !forceDestroy {
				if !confirmAction("Are you sure you want to destroy the entire cluster? This cannot be undone.") {
					fmt.Println("Destroy cancelled.")
					return nil
				}
			}

			tf := platform.NewTerraformExec("./terraform", nebula.ProjectFiles)
			manager := &cluster.Manager{
				TF: tf,
			}

			return manager.Destroy()
		},
	}

	cmd.Flags().BoolVarP(&forceDestroy, "force", "f", false, "Skip confirmation prompts.")
	return cmd
}

func confirmAction(message string) bool {
	var response string
	fmt.Printf("%s (y/n): ", message)
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
