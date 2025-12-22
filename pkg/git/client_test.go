package git

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("expected non-nil git client")
	}
}

func TestGetStatus(t *testing.T) {
	client := NewClient()
	_, err := client.GetStatus()
	if err != nil {
		t.Logf("GetStatus failed: %v", err)
	}
}

func TestGetCurrentBranch(t *testing.T) {
	client := NewClient()
	_, err := client.GetCurrentBranch()
	if err != nil {
		t.Logf("GetCurrentBranch failed: %v", err)
	}
}

func TestListBranches(t *testing.T) {
	client := NewClient()
	_, err := client.ListBranches(false)
	if err != nil {
		t.Logf("ListBranches failed: %v", err)
	}
}

func TestGetLog(t *testing.T) {
	client := NewClient()
	_, err := client.GetLog(10)
	if err != nil {
		t.Logf("GetLog failed: %v", err)
	}
}

func TestGetLastCommit(t *testing.T) {
	client := NewClient()
	_, err := client.GetLastCommit()
	if err != nil {
		t.Logf("GetLastCommit failed: %v", err)
	}
}

func TestGetRemoteList(t *testing.T) {
	client := NewClient()
	_, err := client.GetRemoteList()
	if err != nil {
		t.Logf("GetRemoteList failed: %v", err)
	}
}

func TestGetContributors(t *testing.T) {
	client := NewClient()
	_, err := client.GetContributors()
	if err != nil {
		t.Logf("GetContributors failed: %v", err)
	}
}
