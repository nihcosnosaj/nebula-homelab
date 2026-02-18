package platform

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type AnsibleExec struct {
	Inventory string
	FS        embed.FS
}

func NewAnsibleExec(inventory string, fs embed.FS) *AnsibleExec {
	return &AnsibleExec{Inventory: inventory, FS: fs}
}

// InitializeWorkDir extracts embedded ansible files to a temporary location
func (a *AnsibleExec) InitializeWorkDir() (string, error) {
	tmpDir, err := os.MkdirTemp("", "nebula-ansible-*")
	if err != nil {
		return "", err
	}

	// Walk the embedded "ansible" folder
	err = fs.WalkDir(a.FS, "ansible", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Map "ansible/playbook.yml" -> "/tmp/nebula-ansible-xxx/playbook.yml"
		relPath, _ := filepath.Rel("ansible", path)
		destPath := filepath.Join(tmpDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := fs.ReadFile(a.FS, path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0644)
	})

	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpDir, nil
}

func (a *AnsibleExec) Playbook(playbookPath string) error {
	workDir, err := a.InitializeWorkDir()
	if err != nil {
		return fmt.Errorf("failed to extract ansible files: %w", err)
	}
	defer os.RemoveAll(workDir)

	data, err := os.ReadFile("terraform/inventory.ini")
	if err == nil {
		os.WriteFile(filepath.Join(workDir, "inventory.ini"), data, 0644)
	}

	keyData, err := os.ReadFile("terraform/nebula-key.pem")
	if err == nil {
		os.WriteFile(filepath.Join(workDir, "nebula-key.pem"), keyData, 0400)
	}

	tempInventory := filepath.Join(workDir, "inventory.ini")
	tempPlaybook := filepath.Join(workDir, "cluster-playbook.yml")

	fmt.Println("Verifying SSH connectivity to nodes...")
	if err := a.retryPing(workDir, tempInventory, 5); err != nil {
		return fmt.Errorf("connectivity check failed: %w", err)
	}

	fmt.Println("Running K3s configuration playbook...")
	return a.run(workDir, "ansible-playbook", "-i", tempInventory, tempPlaybook)
}

func (a *AnsibleExec) retryPing(workDir string, inventoryPath string, attempts int) error {
	var err error
	for i := 0; i < attempts; i++ {
		// 'ansible all -m ping' returns 0 if all nodes are reachable
		err = a.run(workDir, "ansible", "all", "-i", inventoryPath, "-m", "ping")
		if err == nil {
			return nil
		}
		fmt.Printf("   (Attempt %d/%d) Nodes not ready yet, retrying in 10s...\n", i+1, attempts)
		time.Sleep(10 * time.Second)
	}
	return err
}

// Helper to run commands within the extracted workDir
func (a *AnsibleExec) run(workDir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "ANSIBLE_FORCE_COLOR=true", "ANSIBLE_HOST_KEY_CHECKING=False")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
