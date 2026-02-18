package platform

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

type TerraformExec struct {
	WorkDir string
	FS      embed.FS
}

func NewTerraformExec(workDir string, fs embed.FS) *TerraformExec {
	return &TerraformExec{WorkDir: workDir, FS: fs}
}

// Apply runs `terraform apply -auto-approve` and streams output to stdout.
func (t *TerraformExec) Apply() error {
	return t.executeInTemp("init", "apply", "-auto-approve")
}

// Destroy runs `terraform destroy -auto-approve`.
func (t *TerraformExec) Destroy() error {
	return t.executeInTemp("init", "destroy", "-auto-approve")
}

// Plan runs `terraform plan`.
func (t *TerraformExec) Plan() error {
	return t.executeInTemp("init", "plan")
}

func (t *TerraformExec) executeInTemp(subcommands ...string) error {
	workDir, err := t.InitializeWorkDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(workDir)

	if err := t.run(workDir, "init"); err != nil {
		return fmt.Errorf("terraform init failed: %w", err)
	}

	if len(subcommands) == 1 && subcommands[0] == "init" {
		return nil
	}

	return t.run(workDir, subcommands[1:]...)
}

func (t *TerraformExec) run(workDir string, args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(),
		"TF_IN_AUTOMATION=true",
		"NO_COLOR=true",
	)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start terraform: %w", err)
	}

	go streamScanner(stdout, "TF [OUT]")
	go streamScanner(stderr, "TF [ERR]")

	return cmd.Wait()
}

func streamScanner(r io.Reader, prefix string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("%s: %s\n", prefix, scanner.Text())
	}
}

func (t *TerraformExec) InitializeWorkDir() (string, error) {
	tmpDir, err := os.MkdirTemp("", "nebula-tf-*")
	if err != nil {
		return "", err
	}

	err = fs.WalkDir(t.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// calculate dest path:
		// "terraform/main.tf" -> "/tmp/nebula-tf-123/main.tf"
		destPath := filepath.Join(tmpDir, path)

		if d.IsDir() {
			// make subdirs in our embed directory
			return os.MkdirAll(destPath, 0755)
		}

		// skip root directory
		if path == "." {
			return nil
		}

		var perm os.FileMode = 0644

		// If it's a provider binary or a shell script, make it executable
		if filepath.Ext(path) == ".sh" || d.Name() == "terraform" || filepath.Base(filepath.Dir(path)) == "darwin_arm64" {
			perm = 0755
		}

		data, err := fs.ReadFile(t.FS, path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, perm)
	})

	if err != nil {
		os.RemoveAll(tmpDir) // clean up if extraction fails
		return "", err
	}

	return filepath.Join(tmpDir, "terraform"), nil
}
