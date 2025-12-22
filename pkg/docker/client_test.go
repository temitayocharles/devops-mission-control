package docker

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("expected non-nil Docker client")
	}
}

func TestListContainers(t *testing.T) {
	client := NewClient()
	_, err := client.ListContainers(true)
	if err != nil {
		t.Logf("docker not available: %v", err)
	}
}

func TestListImages(t *testing.T) {
	client := NewClient()
	_, err := client.ListImages()
	if err != nil {
		t.Logf("docker not available: %v", err)
	}
}

func TestGetSystemInfo(t *testing.T) {
	client := NewClient()
	_, err := client.GetSystemInfo()
	if err != nil {
		t.Logf("docker not available: %v", err)
	}
}

func TestGetContainerStats(t *testing.T) {
	client := NewClient()
	_, err := client.GetContainerStats()
	if err != nil {
		t.Logf("docker not available or no containers: %v", err)
	}
}
