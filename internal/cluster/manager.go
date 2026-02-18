package cluster

import (
	"fmt"
	"strings"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nihcosnosaj/nebula-homelab/internal/platform"
)

type Manager struct {
	TF      TerraformProvider
	Ansible AnsibleProvider
}

type TerraformProvider interface {
	Apply() error
	Destroy() error
	Plan() error
}

type AnsibleProvider interface {
	Playbook(path string) error
}

// Up handles orchestration.
func (m *Manager) Up(dryRun bool) error {
	if dryRun {
		fmt.Println("Performing a dryrun (Terraform Plan)...")
		return m.TF.Plan()
	}

	fmt.Println("Provisioning infrastructure...")
	if err := m.TF.Apply(); err != nil {
		return fmt.Errorf("terraform apply failed: %w", err)
	}

	fmt.Println("Configuring control-plane and worker nodes...")
	if err := m.Ansible.Playbook("ansible/cluster-playbook.yml"); err != nil {
		return fmt.Errorf("ansible-playbook failed: %w", err)
	}

	return nil
}

func (m *Manager) Destroy() error {
	fmt.Println("Destroying cluster and bringing down AWS instances...")

	// cleanup playbooks for stateful data

	if err := m.TF.Destroy(); err != nil {
		return fmt.Errorf("infrastructure destruction failed: %w", err)
	}

	fmt.Println("All resources successfully destroyed.")
	return nil
}

func (m *Manager) Status(projectName string) error {
	awsSvc, _ := platform.NewAWSClient("us-west-1")
	instances, err := awsSvc.GetClusterInstance(projectName)
	if err != nil {
		return fmt.Errorf("failed getting cluster instance: %w", err)
	}

	var totalBurn float64
	fmt.Printf("Nebula Cluster: %s\n", projectName)
	fmt.Printf("%-20s %-12s %-12s %-12s %-10s %-10s\n", "NAME", "TYPE", "AZ", "LIVE PRICE", "STATE", "HEALTH")
	fmt.Println(strings.Repeat("-", 80))

	for _, inst := range instances {
		name := getTagName(inst.Tags, "Name")

		interrupted, _, _ := awsSvc.CheckForInterruptions(*inst.InstanceId)
		healthStatus := "Healthy"
		if interrupted {
			healthStatus = "Reclaimed"
		}
		price, _ := awsSvc.GetCurrentSpotPrice(inst.InstanceType, *inst.Placement.AvailabilityZone)
		totalBurn += price

		fmt.Printf("%-20s %-12s %-12s $%0.4f/hr  %-10s\n %-10s\n",
			name, inst.InstanceType, *inst.Placement.AvailabilityZone, price, string(inst.State.Name), healthStatus)
	}

	fmt.Println(strings.Repeat("-", 75))
	fmt.Printf("Total Estimated Burn: $%0.4f/hr (~$%0.2f/month)\n\n", totalBurn, totalBurn*24*30)
	return nil
}

func getTagName(tags []ec2types.Tag, key string) string {
	for _, t := range tags {
		if *t.Key == key {
			return *t.Value
		}
	}
	return "unknown"
}
