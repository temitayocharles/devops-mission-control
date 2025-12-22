package docker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Client wraps Docker operations
type Client struct{}

// NewClient creates a new Docker client
func NewClient() *Client {
	return &Client{}
}

// execDocker runs a docker command and returns output
func (c *Client) execDocker(args ...string) (string, error) {
	cmd := exec.Command("docker", args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// ListContainers lists all containers (running and stopped)
func (c *Client) ListContainers(running bool) (string, error) {
	args := []string{"ps", "-a", "--format", "table {{.Names}}\t{{.Status}}\t{{.Ports}}"}
	if running {
		args = []string{"ps", "--format", "table {{.Names}}\t{{.Status}}\t{{.Ports}}"}
	}
	return c.execDocker(args...)
}

// ListImages lists all Docker images
func (c *Client) ListImages() (string, error) {
	return c.execDocker("images", "--format", "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}")
}

// StopContainer stops a running container
func (c *Client) StopContainer(containerID string) (string, error) {
	return c.execDocker("stop", containerID)
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(containerID string, force bool) (string, error) {
	args := []string{"rm", containerID}
	if force {
		args = append(args, "-f")
	}
	return c.execDocker(args...)
}

// GetContainerLogs retrieves logs from a container
func (c *Client) GetContainerLogs(containerID string, follow bool) error {
	args := []string{"logs", containerID, "--tail", "50"}
	if follow {
		args = append(args, "-f")
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ExecInContainer executes a command in a container (interactive)
func (c *Client) ExecInContainer(containerID string, command []string) error {
	args := []string{"exec", "-it", containerID}
	args = append(args, command...)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// GetContainerStats retrieves container resource usage stats
func (c *Client) GetContainerStats() (string, error) {
	return c.execDocker("stats", "--no-stream", "--format", "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}")
}

// GetContainerInspect retrieves detailed info about a container
func (c *Client) GetContainerInspect(containerID string) (string, error) {
	return c.execDocker("inspect", containerID)
}

// SystemPrune removes unused Docker resources
func (c *Client) SystemPrune(all bool) (string, error) {
	args := []string{"system", "prune", "-f"}
	if all {
		args = append(args, "-a")
	}
	return c.execDocker(args...)
}

// GetSystemInfo retrieves Docker system information
func (c *Client) GetSystemInfo() (string, error) {
	return c.execDocker("system", "df")
}

// ComposeUp starts Docker Compose services
func (c *Client) ComposeUp(detach bool) error {
	args := []string{"compose", "up"}
	if detach {
		args = append(args, "-d")
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ComposeDown stops Docker Compose services
func (c *Client) ComposeDown() error {
	cmd := exec.Command("docker", "compose", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ComposeLogs streams Docker Compose logs
func (c *Client) ComposeLogs(follow bool) error {
	args := []string{"compose", "logs"}
	if follow {
		args = append(args, "-f")
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ComposeStatus shows Docker Compose service status
func (c *Client) ComposeStatus() (string, error) {
	return c.execDocker("compose", "ps")
}

// BuildImage builds a Docker image
func (c *Client) BuildImage(dockerfile, tag, context string) error {
	args := []string{"build", "-f", dockerfile, "-t", tag}
	if context != "" {
		args = append(args, context)
	} else {
		args = append(args, ".")
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// PullImage pulls an image from registry
func (c *Client) PullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// RemoveImage removes a Docker image
func (c *Client) RemoveImage(imageID string, force bool) (string, error) {
	args := []string{"rmi", imageID}
	if force {
		args = append(args, "-f")
	}
	return c.execDocker(args...)
}
