package azure

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps Azure operations
type Client struct {
	Subscription  string
	ResourceGroup string
}

// NewClient creates a new Azure client
func NewClient(subscription, resourceGroup string) *Client {
	if subscription == "" {
		subscription = getDefaultSubscription()
	}
	return &Client{
		Subscription:  subscription,
		ResourceGroup: resourceGroup,
	}
}

// execAZ runs an az command and returns output
func (c *Client) execAZ(args ...string) (string, error) {
	fullArgs := append([]string{"--subscription=" + c.Subscription}, args...)

	cmd := exec.Command("az", fullArgs...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("az error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// getDefaultSubscription gets the default Azure subscription
func getDefaultSubscription() string {
	cmd := exec.Command("az", "account", "show", "--query=id", "-o=tsv")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}

// VM operations
func (c *Client) ListVMs() (string, error) {
	args := []string{"vm", "list"}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(append(args, "-d")...)
}

func (c *Client) StartVM(name string) (string, error) {
	args := []string{"vm", "start", "--name=" + name}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

func (c *Client) StopVM(name string) (string, error) {
	args := []string{"vm", "deallocate", "--name=" + name}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

func (c *Client) DeleteVM(name string) (string, error) {
	args := []string{"vm", "delete", "--name=" + name, "--yes"}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

// Storage Account operations
func (c *Client) ListStorageAccounts() (string, error) {
	args := []string{"storage", "account", "list"}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

func (c *Client) ListStorageContainers(accountName string) (string, error) {
	return c.execAZ("storage", "container", "list", "--account-name="+accountName)
}

func (c *Client) ListStorageBlobs(accountName, container string) (string, error) {
	return c.execAZ("storage", "blob", "list", "--account-name="+accountName, "--container-name="+container)
}

func (c *Client) GetStorageAccountKeys(name string) (string, error) {
	args := []string{"storage", "account", "keys", "list", "--account-name=" + name}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

// Database operations
func (c *Client) ListDatabases() (string, error) {
	args := []string{"sql", "db", "list"}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

func (c *Client) ListSQLServers() (string, error) {
	args := []string{"sql", "server", "list"}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

// App Service operations
func (c *Client) ListAppServices() (string, error) {
	args := []string{"appservice", "list"}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

func (c *Client) GetAppServiceConfig(name string) (string, error) {
	args := []string{"appservice", "config", "show", "--name=" + name}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

// Resource Group operations
func (c *Client) ListResourceGroups() (string, error) {
	return c.execAZ("group", "list")
}

func (c *Client) GetResourceGroupInfo(name string) (string, error) {
	return c.execAZ("group", "show", "--name="+name)
}

// Network operations
func (c *Client) ListNetworks() (string, error) {
	args := []string{"network", "vnet", "list"}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

func (c *Client) ListNetworkSecurityGroups() (string, error) {
	args := []string{"network", "nsg", "list"}
	if c.ResourceGroup != "" {
		args = append(args, "--resource-group="+c.ResourceGroup)
	}
	return c.execAZ(args...)
}

// Account and subscription
func (c *Client) GetAccountInfo() (string, error) {
	return c.execAZ("account", "show")
}

func (c *Client) ListSubscriptions() (string, error) {
	return c.execAZ("account", "list")
}

func (c *Client) GetCostManagementData() (string, error) {
	return c.execAZ("costmanagement", "query", "usage")
}
