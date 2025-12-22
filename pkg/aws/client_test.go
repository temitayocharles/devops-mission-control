package aws

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("us-east-1", "default")
	if client.Region != "us-east-1" {
		t.Errorf("expected region 'us-east-1', got %s", client.Region)
	}
	if client.Profile != "default" {
		t.Errorf("expected profile 'default', got %s", client.Profile)
	}
}

func TestNewClientEmpty(t *testing.T) {
	client := NewClient("", "")
	if client == nil {
		t.Error("expected non-nil AWS client")
	}
}

func TestGetCallerIdentity(t *testing.T) {
	client := NewClient("", "")
	_, err := client.GetCallerIdentity()
	if err != nil {
		t.Logf("aws cli not available or not configured: %v", err)
	}
}

func TestListEC2Instances(t *testing.T) {
	client := NewClient("us-east-1", "")
	_, err := client.ListEC2Instances(false)
	if err != nil {
		t.Logf("aws cli not available: %v", err)
	}
}

func TestListS3Buckets(t *testing.T) {
	client := NewClient("us-east-1", "")
	_, err := client.ListS3Buckets()
	if err != nil {
		t.Logf("aws cli not available: %v", err)
	}
}

func TestListRDSInstances(t *testing.T) {
	client := NewClient("us-east-1", "")
	_, err := client.ListRDSInstances()
	if err != nil {
		t.Logf("aws cli not available: %v", err)
	}
}

func TestFindUnusedResources(t *testing.T) {
	client := NewClient("us-east-1", "")
	resources, err := client.FindUnusedResources()
	if err != nil {
		t.Logf("aws cli not available: %v", err)
	}
	if resources == nil {
		t.Error("resources map should not be nil")
	}
}

func TestListSecurityGroups(t *testing.T) {
	client := NewClient("us-east-1", "")
	_, err := client.ListSecurityGroups()
	if err != nil {
		t.Logf("aws cli not available: %v", err)
	}
}

func TestListVPCs(t *testing.T) {
	client := NewClient("us-east-1", "")
	_, err := client.ListVPCs()
	if err != nil {
		t.Logf("aws cli not available: %v", err)
	}
}
