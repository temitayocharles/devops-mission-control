package k8s

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Client wraps kubectl operations
type Client struct {
	Namespace string
	Context   string
}

// NewClient creates a new Kubernetes client
func NewClient(namespace, context string) *Client {
	if namespace == "" {
		namespace = "default"
	}
	return &Client{
		Namespace: namespace,
		Context:   context,
	}
}

// execKubectl runs a kubectl command and returns output
func (c *Client) execKubectl(args ...string) (string, error) {
	if c.Context != "" {
		args = append([]string{"--context", c.Context}, args...)
	}

	cmd := exec.Command("kubectl", args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("kubectl error: %s", errOut.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// ListPods lists pods in a namespace
func (c *Client) ListPods(namespace string) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("get", "pods", "-n", namespace, "-o", "wide")
}

// GetPodLogs retrieves logs from a pod
func (c *Client) GetPodLogs(podName, namespace string, follow bool) error {
	if namespace == "" {
		namespace = c.Namespace
	}
	args := []string{"logs", podName, "-n", namespace, "--tail=50"}
	if follow {
		args = append(args, "-f")
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ListDeployments lists deployments in a namespace
func (c *Client) ListDeployments(namespace string) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("get", "deployments", "-n", namespace, "-o", "wide")
}

// ListServices lists services in a namespace
func (c *Client) ListServices(namespace string) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("get", "services", "-n", namespace, "-o", "wide")
}

// ListNodes lists all cluster nodes
func (c *Client) ListNodes() (string, error) {
	return c.execKubectl("get", "nodes", "-o", "wide")
}

// GetCurrentContext returns the current cluster context
func (c *Client) GetCurrentContext() (string, error) {
	return c.execKubectl("config", "current-context")
}

// ListContexts lists all available contexts
func (c *Client) ListContexts() (string, error) {
	return c.execKubectl("config", "get-contexts")
}

// SwitchContext switches to a different context
func (c *Client) SwitchContext(contextName string) (string, error) {
	return c.execKubectl("config", "use-context", contextName)
}

// DescribePod describes a pod
func (c *Client) DescribePod(podName, namespace string) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("describe", "pod", podName, "-n", namespace)
}

// DescribeDeployment describes a deployment
func (c *Client) DescribeDeployment(deploymentName, namespace string) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("describe", "deployment", deploymentName, "-n", namespace)
}

// DeletePod deletes a pod
func (c *Client) DeletePod(podName, namespace string) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("delete", "pod", podName, "-n", namespace)
}

// ScaleDeployment scales a deployment
func (c *Client) ScaleDeployment(deploymentName, namespace string, replicas int) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("scale", "deployment", deploymentName, "-n", namespace, "--replicas", fmt.Sprintf("%d", replicas))
}

// GetEvents gets cluster events
func (c *Client) GetEvents(namespace string) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("get", "events", "-n", namespace, "--sort-by='.lastTimestamp'")
}

// GetPodsByLabel gets pods matching a label selector
func (c *Client) GetPodsByLabel(selector, namespace string) (string, error) {
	if namespace == "" {
		namespace = c.Namespace
	}
	return c.execKubectl("get", "pods", "-n", namespace, "-l", selector, "-o", "wide")
}

// ExecInPod executes a command in a pod (interactive)
func (c *Client) ExecInPod(podName, namespace, container string, command []string) error {
	if namespace == "" {
		namespace = c.Namespace
	}

	args := []string{"exec", "-it", podName, "-n", namespace}
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, "--")
	args = append(args, command...)

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// PortForward forwards a local port to a pod
func (c *Client) PortForward(podName, namespace, ports string) error {
	if namespace == "" {
		namespace = c.Namespace
	}

	cmd := exec.Command("kubectl", "port-forward", podName, ports, "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// GetClusterInfo gets cluster information
func (c *Client) GetClusterInfo() (string, error) {
	return c.execKubectl("cluster-info")
}

// CheckClusterHealth performs basic health checks
func (c *Client) CheckClusterHealth() (map[string]string, error) {
	health := make(map[string]string)

	// Check nodes
	nodes, err := c.execKubectl("get", "nodes", "-o", "jsonpath={.items[*].metadata.name}")
	if err != nil {
		health["nodes"] = "ERROR: " + err.Error()
	} else {
		health["nodes"] = fmt.Sprintf("OK (%s)", nodes)
	}

	// Check API server
	_, err1 := c.execKubectl("get", "componentstatuses")
	if err1 == nil {
		health["api"] = "OK"
	} else {
		health["api"] = "ERROR"
	}

	// Check namespaces
	_, err2 := c.execKubectl("get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	if err2 != nil {
		health["namespaces"] = "ERROR"
	} else {
		health["namespaces"] = "OK"
	}

	return health, nil
}
