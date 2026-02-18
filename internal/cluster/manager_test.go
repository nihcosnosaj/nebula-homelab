package cluster

import (
	"errors"
	"testing"
)

type MockTerraform struct {
	ApplyCalled bool
	ShouldFail  bool
}

func (m *MockTerraform) Apply() error {
	m.ApplyCalled = true
	if m.ShouldFail {
		return errors.New("tf error")
	}
	return nil
}

func (m *MockTerraform) Destroy() error { return nil }
func (m *MockTerraform) Plan() error    { return nil }

type MockAnsible struct {
	PlaybookCalled bool
}

func (m *MockAnsible) Playbook(path string) error {
	m.PlaybookCalled = true
	return nil
}

func TestManager_Up(t *testing.T) {
	t.Run("Success Path", func(t *testing.T) {
		tf := &MockTerraform{}
		ang := &MockAnsible{}
		mgr := &Manager{TF: tf, Ansible: ans}

		err := mgr.Up(false) // no dry-run

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !tf.ApplyCalled || !ans.PlaybookCalled {
			t.Error("expected both TF and Ansible to be called")
		}
	})

	t.Run("Terraform Fails", func(t *testing.T) {
		tf := &MockTerraform{ShouldFail: true}
		ans := &MockAnsible{}
		mgr := &Manager{TF: tf, Ansible: ans}

		err := mgr.Up(false)

		if err == nil {
			t.Error("expected error from failed terraform, got nil")
		}
		if ans.PlaybookCalled {
			t.Error("ansible shouldn't run if terraform fails...")
		}
	})
}
