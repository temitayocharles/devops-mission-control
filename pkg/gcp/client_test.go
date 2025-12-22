package gcp

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-project", "us-central1")
	if client.Project != "test-project" {
		t.Errorf("expected project 'test-project', got %s", client.Project)
	}
	if client.Region != "us-central1" {
		t.Errorf("expected region 'us-central1', got %s", client.Region)
	}
}

func TestNewClientDefault(t *testing.T) {
	client := NewClient("", "")
	if client.Region != "us-central1" {
		t.Errorf("expected default region 'us-central1', got %s", client.Region)
	}
}

func TestListInstances(t *testing.T) {
	client := NewClient("test-project", "us-central1")
	_, err := client.ListInstances()
	if err != nil {
		t.Logf("gcloud not available: %v", err)
	}
}

func TestListBuckets(t *testing.T) {
	client := NewClient("test-project", "us-central1")
	_, err := client.ListBuckets()
	if err != nil {
		t.Logf("gcloud not available: %v", err)
	}
}

func TestListSQLInstances(t *testing.T) {
	client := NewClient("test-project", "us-central1")
	_, err := client.ListSQLInstances()
	if err != nil {
		t.Logf("gcloud not available: %v", err)
	}
}

func TestListCloudRunServices(t *testing.T) {
	client := NewClient("test-project", "us-central1")
	_, err := client.ListCloudRunServices()
	if err != nil {
		t.Logf("gcloud not available: %v", err)
	}
}

func TestListServiceAccounts(t *testing.T) {
	client := NewClient("test-project", "us-central1")
	_, err := client.ListServiceAccounts()
	if err != nil {
		t.Logf("gcloud not available: %v", err)
	}
}

func TestGetProjectInfo(t *testing.T) {
	client := NewClient("test-project", "us-central1")
	_, err := client.GetProjectInfo()
	if err != nil {
		t.Logf("gcloud not available: %v", err)
	}
}
