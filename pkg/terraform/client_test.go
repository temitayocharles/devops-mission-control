package terraform

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("/tmp/tf-test")
	if client.WorkDir != "/tmp/tf-test" {
		t.Errorf("expected workdir '/tmp/tf-test', got %s", client.WorkDir)
	}
}

func TestNewClientDefault(t *testing.T) {
	client := NewClient("")
	if client.WorkDir != "." {
		t.Errorf("expected default workdir '.', got %s", client.WorkDir)
	}
}

func TestVersion(t *testing.T) {
	client := NewClient(".")
	_, err := client.Version()
	if err != nil {
		t.Logf("terraform not available: %v", err)
	}
}

func TestValidate(t *testing.T) {
	client := NewClient(".")
	_, err := client.Validate()
	if err != nil {
		t.Logf("terraform not available or no terraform files: %v", err)
	}
}

func TestFormat(t *testing.T) {
	client := NewClient(".")
	_, err := client.Format()
	if err != nil {
		t.Logf("terraform not available: %v", err)
	}
}

func TestWorkspaceList(t *testing.T) {
	client := NewClient(".")
	_, err := client.WorkspaceList()
	if err != nil {
		t.Logf("terraform not available: %v", err)
	}
}

func TestStateList(t *testing.T) {
	client := NewClient(".")
	_, err := client.StateList()
	if err != nil {
		t.Logf("terraform or terraform state not available: %v", err)
	}
}

func TestOutputAll(t *testing.T) {
	client := NewClient(".")
	_, err := client.OutputAll()
	if err != nil {
		t.Logf("terraform not available: %v", err)
	}
}
