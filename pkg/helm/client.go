package helm

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps Helm operations
type Client struct {
	Namespace string
}

// NewClient creates a new Helm client
func NewClient(namespace string) *Client {
	if namespace == "" {
		namespace = "default"
	}
	return &Client{
		Namespace: namespace,
	}
}

// execHelm runs a helm command and returns output
func (c *Client) execHelm(args ...string) (string, error) {
	fullArgs := append([]string{"-n", c.Namespace}, args...)

	cmd := exec.Command("helm", fullArgs...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("helm error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Repository operations
func (c *Client) RepoList() (string, error) {
	return c.execHelm("repo", "list")
}

func (c *Client) RepoAdd(name, url string) (string, error) {
	return c.execHelm("repo", "add", name, url)
}

func (c *Client) RepoRemove(name string) (string, error) {
	return c.execHelm("repo", "remove", name)
}

func (c *Client) RepoUpdate() (string, error) {
	return c.execHelm("repo", "update")
}

// Chart operations
func (c *Client) SearchChart(name string) (string, error) {
	return c.execHelm("search", "hub", name)
}

func (c *Client) ShowChart(chart string) (string, error) {
	return c.execHelm("show", "chart", chart)
}

func (c *Client) ShowValues(chart string) (string, error) {
	return c.execHelm("show", "values", chart)
}

// Release operations
func (c *Client) ListReleases() (string, error) {
	return c.execHelm("list")
}

func (c *Client) GetRelease(name string) (string, error) {
	return c.execHelm("get", "all", name)
}

func (c *Client) GetReleaseValues(name string) (string, error) {
	return c.execHelm("get", "values", name)
}

func (c *Client) GetReleaseStatus(name string) (string, error) {
	return c.execHelm("status", name)
}

func (c *Client) InstallChart(name, chart string, values map[string]interface{}) (string, error) {
	args := []string{"install", name, chart}
	for k, v := range values {
		args = append(args, "--set", fmt.Sprintf("%s=%v", k, v))
	}
	return c.execHelm(args...)
}

func (c *Client) UpgradeRelease(name, chart string) (string, error) {
	return c.execHelm("upgrade", name, chart)
}

func (c *Client) UninstallRelease(name string) (string, error) {
	return c.execHelm("uninstall", name)
}

func (c *Client) RollbackRelease(name string, revision int) (string, error) {
	return c.execHelm("rollback", name, fmt.Sprintf("%d", revision))
}

// Release history
func (c *Client) GetReleaseHistory(name string) (string, error) {
	return c.execHelm("history", name)
}

// Template operations
func (c *Client) TemplateChart(name, chart string) (string, error) {
	return c.execHelm("template", name, chart)
}

// Lint chart
func (c *Client) LintChart(path string) (string, error) {
	return c.execHelm("lint", path)
}

// Plugin operations
func (c *Client) ListPlugins() (string, error) {
	return c.execHelm("plugin", "list")
}

// Version
func (c *Client) GetVersion() (string, error) {
	return c.execHelm("version")
}
