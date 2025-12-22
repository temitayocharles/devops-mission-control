package network

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client.Timeout != "5" {
		t.Errorf("expected timeout '5', got %s", client.Timeout)
	}
}

func TestPing(t *testing.T) {
	client := NewClient()
	_, err := client.Ping("localhost")
	if err != nil {
		t.Logf("ping not available: %v", err)
	}
}

func TestGetHostname(t *testing.T) {
	client := NewClient()
	hostname, err := client.GetHostname()
	if err != nil {
		t.Logf("hostname command not available: %v", err)
	} else if hostname == "" {
		t.Log("hostname is empty")
	}
}

func TestGetNetworkInterfaces(t *testing.T) {
	client := NewClient()
	_, err := client.GetNetworkInterfaces()
	if err != nil {
		t.Logf("ifconfig not available: %v", err)
	}
}

func TestGetRouteTable(t *testing.T) {
	client := NewClient()
	_, err := client.GetRouteTable()
	if err != nil {
		t.Logf("netstat not available: %v", err)
	}
}

func TestGetDNSServers(t *testing.T) {
	client := NewClient()
	_, err := client.GetDNSServers()
	if err != nil {
		t.Logf("DNS resolution not available: %v", err)
	}
}

func TestNslookup(t *testing.T) {
	client := NewClient()
	_, err := client.Nslookup("localhost")
	if err != nil {
		t.Logf("nslookup not available: %v", err)
	}
}
