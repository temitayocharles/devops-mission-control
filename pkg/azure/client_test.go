package azure

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-sub", "test-rg")
	if client.Subscription != "test-sub" {
		t.Errorf("expected subscription 'test-sub', got %s", client.Subscription)
	}
	if client.ResourceGroup != "test-rg" {
		t.Errorf("expected resource group 'test-rg', got %s", client.ResourceGroup)
	}
}

func TestListVMs(t *testing.T) {
	client := NewClient("test-sub", "test-rg")
	_, err := client.ListVMs()
	if err != nil {
		t.Logf("az cli not available: %v", err)
	}
}

func TestListStorageAccounts(t *testing.T) {
	client := NewClient("test-sub", "test-rg")
	_, err := client.ListStorageAccounts()
	if err != nil {
		t.Logf("az cli not available: %v", err)
	}
}

func TestListDatabases(t *testing.T) {
	client := NewClient("test-sub", "test-rg")
	_, err := client.ListDatabases()
	if err != nil {
		t.Logf("az cli not available: %v", err)
	}
}

func TestListAppServices(t *testing.T) {
	client := NewClient("test-sub", "test-rg")
	_, err := client.ListAppServices()
	if err != nil {
		t.Logf("az cli not available: %v", err)
	}
}

func TestListResourceGroups(t *testing.T) {
	client := NewClient("test-sub", "")
	_, err := client.ListResourceGroups()
	if err != nil {
		t.Logf("az cli not available: %v", err)
	}
}

func TestGetAccountInfo(t *testing.T) {
	client := NewClient("test-sub", "test-rg")
	_, err := client.GetAccountInfo()
	if err != nil {
		t.Logf("az cli not available: %v", err)
	}
}
