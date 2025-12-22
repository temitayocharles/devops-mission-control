package terraform

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps Terraform operations
type Client struct {
	WorkDir string
}

// NewClient creates a new Terraform client
func NewClient(workDir string) *Client {
	if workDir == "" {
		workDir = "."
	}
	return &Client{
		WorkDir: workDir,
	}
}

// execTerraform runs a terraform command in the working directory
func (c *Client) execTerraform(args ...string) (string, error) {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = c.WorkDir
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("terraform error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Version
func (c *Client) Version() (string, error) {
	return c.execTerraform("version")
}

// Init
func (c *Client) Init() (string, error) {
	return c.execTerraform("init")
}

// Validate configuration
func (c *Client) Validate() (string, error) {
	return c.execTerraform("validate")
}

// Format configuration
func (c *Client) Format() (string, error) {
	return c.execTerraform("fmt", "-recursive")
}

// Plan
func (c *Client) Plan(varFile string) (string, error) {
	args := []string{"plan"}
	if varFile != "" {
		args = append(args, "-var-file="+varFile)
	}
	return c.execTerraform(args...)
}

// Plan with output file
func (c *Client) PlanToFile(varFile, outFile string) (string, error) {
	args := []string{"plan", "-out=" + outFile}
	if varFile != "" {
		args = append(args, "-var-file="+varFile)
	}
	return c.execTerraform(args...)
}

// Apply
func (c *Client) Apply(planFile string) (string, error) {
	args := []string{"apply"}
	if planFile != "" {
		args = append(args, planFile)
	} else {
		args = append(args, "-auto-approve")
	}
	return c.execTerraform(args...)
}

// Destroy
func (c *Client) Destroy(varFile string) (string, error) {
	args := []string{"destroy", "-auto-approve"}
	if varFile != "" {
		args = append(args, "-var-file="+varFile)
	}
	return c.execTerraform(args...)
}

// State operations
func (c *Client) StateList() (string, error) {
	return c.execTerraform("state", "list")
}

func (c *Client) StateShow(resource string) (string, error) {
	return c.execTerraform("state", "show", resource)
}

func (c *Client) StateRm(resource string) (string, error) {
	return c.execTerraform("state", "rm", resource)
}

func (c *Client) StateMv(oldResource, newResource string) (string, error) {
	return c.execTerraform("state", "mv", oldResource, newResource)
}

// Output
func (c *Client) Output(name string) (string, error) {
	return c.execTerraform("output", name)
}

func (c *Client) OutputAll() (string, error) {
	return c.execTerraform("output", "-json")
}

// Modules
func (c *Client) GetModules() (string, error) {
	return c.execTerraform("get", "-update")
}

// Console
func (c *Client) Console() (string, error) {
	return c.execTerraform("console")
}

// Import resource
func (c *Client) Import(address, id string) (string, error) {
	return c.execTerraform("import", address, id)
}

// Taint resource
func (c *Client) Taint(address string) (string, error) {
	return c.execTerraform("taint", address)
}

// Untaint resource
func (c *Client) Untaint(address string) (string, error) {
	return c.execTerraform("untaint", address)
}

// Workspace operations
func (c *Client) WorkspaceList() (string, error) {
	return c.execTerraform("workspace", "list")
}

func (c *Client) WorkspaceNew(name string) (string, error) {
	return c.execTerraform("workspace", "new", name)
}

func (c *Client) WorkspaceSelect(name string) (string, error) {
	return c.execTerraform("workspace", "select", name)
}

func (c *Client) WorkspaceDelete(name string) (string, error) {
	return c.execTerraform("workspace", "delete", name)
}
