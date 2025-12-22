package helm

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-ns")
	if client.Namespace != "test-ns" {
		t.Errorf("expected namespace 'test-ns', got %s", client.Namespace)
	}
}

func TestNewClientDefault(t *testing.T) {
	client := NewClient("")
	if client.Namespace != "default" {
		t.Errorf("expected default namespace, got %s", client.Namespace)
	}
}

func TestListReleases(t *testing.T) {
	client := NewClient("default")
	_, err := client.ListReleases()
	if err != nil {
		t.Logf("helm not available: %v", err)
	}
}

func TestRepoList(t *testing.T) {
	client := NewClient("default")
	_, err := client.RepoList()
	if err != nil {
		t.Logf("helm not available: %v", err)
	}
}

func TestSearchChart(t *testing.T) {
	client := NewClient("default")
	_, err := client.SearchChart("nginx")
	if err != nil {
		t.Logf("helm not available: %v", err)
	}
}

func TestGetVersion(t *testing.T) {
	client := NewClient("default")
	_, err := client.GetVersion()
	if err != nil {
		t.Logf("helm not available: %v", err)
	}
}

func TestListPlugins(t *testing.T) {
	client := NewClient("default")
	_, err := client.ListPlugins()
	if err != nil {
		t.Logf("helm not available: %v", err)
	}
}
