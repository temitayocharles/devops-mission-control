package k8s

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-ns", "test-context")
	if client.Namespace != "test-ns" {
		t.Errorf("expected namespace 'test-ns', got %s", client.Namespace)
	}
	if client.Context != "test-context" {
		t.Errorf("expected context 'test-context', got %s", client.Context)
	}
}

func TestNewClientDefault(t *testing.T) {
	client := NewClient("", "")
	if client.Namespace != "default" {
		t.Errorf("expected default namespace, got %s", client.Namespace)
	}
}

func TestGetCurrentContext(t *testing.T) {
	client := NewClient("", "")
	// This test requires kubectl to be installed and configured
	// For CI/CD, we'd mock this, but for dev purposes this is useful
	_, err := client.GetCurrentContext()
	if err != nil {
		t.Logf("kubectl not available or not configured: %v", err)
	}
}

func TestListPods(t *testing.T) {
	client := NewClient("default", "")
	_, err := client.ListPods("default")
	if err != nil {
		t.Logf("kubectl not available: %v", err)
	}
}

func TestListDeployments(t *testing.T) {
	client := NewClient("default", "")
	_, err := client.ListDeployments("default")
	if err != nil {
		t.Logf("kubectl not available: %v", err)
	}
}

func TestListServices(t *testing.T) {
	client := NewClient("default", "")
	_, err := client.ListServices("default")
	if err != nil {
		t.Logf("kubectl not available: %v", err)
	}
}

func TestListNodes(t *testing.T) {
	client := NewClient("", "")
	_, err := client.ListNodes()
	if err != nil {
		t.Logf("kubectl not available: %v", err)
	}
}

func TestCheckClusterHealth(t *testing.T) {
	client := NewClient("", "")
	health, err := client.CheckClusterHealth()
	if err != nil {
		t.Logf("health check failed (kubectl may not be configured): %v", err)
	}
	if health == nil {
		t.Error("health map should not be nil")
	}
}
