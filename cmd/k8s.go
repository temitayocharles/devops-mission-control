package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/ops-tool/pkg/k8s"
)

var k8sNamespace string
var k8sContext string

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes operations",
	Long:  "Manage Kubernetes clusters, pods, deployments, and services",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize any k8s-specific setup here
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var k8sPodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Manage Kubernetes pods",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var k8sPodsListCmd = &cobra.Command{
	Use:   "list [namespace]",
	Short: "List pods in a namespace",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns := k8sNamespace
		if len(args) > 0 {
			ns = args[0]
		}

		client := k8s.NewClient(ns, k8sContext)
		output, err := client.ListPods(ns)
		if err != nil {
			return fmt.Errorf("failed to list pods: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var k8sPodsLogsCmd = &cobra.Command{
	Use:   "logs <pod-name> [namespace]",
	Short: "Stream logs from a pod",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		podName := args[0]
		ns := k8sNamespace
		if len(args) > 1 {
			ns = args[1]
		}

		follow, _ := cmd.Flags().GetBool("follow")
		client := k8s.NewClient(ns, k8sContext)
		return client.GetPodLogs(podName, ns, follow)
	},
}

var k8sPodsDescribeCmd = &cobra.Command{
	Use:   "describe <pod-name> [namespace]",
	Short: "Describe a pod",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		podName := args[0]
		ns := k8sNamespace
		if len(args) > 1 {
			ns = args[1]
		}

		client := k8s.NewClient(ns, k8sContext)
		output, err := client.DescribePod(podName, ns)
		if err != nil {
			return fmt.Errorf("failed to describe pod: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var k8sPodsDeleteCmd = &cobra.Command{
	Use:   "delete <pod-name> [namespace]",
	Short: "Delete a pod",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		podName := args[0]
		ns := k8sNamespace
		if len(args) > 1 {
			ns = args[1]
		}

		client := k8s.NewClient(ns, k8sContext)
		output, err := client.DeletePod(podName, ns)
		if err != nil {
			return fmt.Errorf("failed to delete pod: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var k8sDeploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "Manage Kubernetes deployments",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var k8sDeploymentsListCmd = &cobra.Command{
	Use:   "list [namespace]",
	Short: "List deployments",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns := k8sNamespace
		if len(args) > 0 {
			ns = args[0]
		}

		client := k8s.NewClient(ns, k8sContext)
		output, err := client.ListDeployments(ns)
		if err != nil {
			return fmt.Errorf("failed to list deployments: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var k8sDeploymentsScaleCmd = &cobra.Command{
	Use:   "scale <deployment-name> <replicas> [namespace]",
	Short: "Scale a deployment",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		deploymentName := args[0]
		replicas, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("replicas must be a number")
		}

		ns := k8sNamespace
		if len(args) > 2 {
			ns = args[2]
		}

		client := k8s.NewClient(ns, k8sContext)
		output, err := client.ScaleDeployment(deploymentName, ns, replicas)
		if err != nil {
			return fmt.Errorf("failed to scale deployment: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var k8sServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Manage Kubernetes services",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var k8sServicesListCmd = &cobra.Command{
	Use:   "list [namespace]",
	Short: "List services",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ns := k8sNamespace
		if len(args) > 0 {
			ns = args[0]
		}

		client := k8s.NewClient(ns, k8sContext)
		output, err := client.ListServices(ns)
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var k8sNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage Kubernetes nodes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var k8sNodesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := k8s.NewClient("", k8sContext)
		output, err := client.ListNodes()
		if err != nil {
			return fmt.Errorf("failed to list nodes: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var k8sContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage Kubernetes contexts",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

var k8sContextListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available contexts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := k8s.NewClient("", "")
		output, err := client.ListContexts()
		if err != nil {
			return fmt.Errorf("failed to list contexts: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var k8sContextCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current context",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := k8s.NewClient("", "")
		context, err := client.GetCurrentContext()
		if err != nil {
			return fmt.Errorf("failed to get current context: %w", err)
		}

		fmt.Printf("Current context: %s\n", context)
		return nil
	},
}

var k8sContextSwitchCmd = &cobra.Command{
	Use:   "switch <context>",
	Short: "Switch to a different context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contextName := args[0]
		client := k8s.NewClient("", "")
		_, err := client.SwitchContext(contextName)
		if err != nil {
			return fmt.Errorf("failed to switch context: %w", err)
		}

		fmt.Printf("✅ Switched to context: %s\n", contextName)
		return nil
	},
}

var k8sHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check cluster health",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := k8s.NewClient("", k8sContext)
		health, err := client.CheckClusterHealth()
		if err != nil {
			return fmt.Errorf("failed to check health: %w", err)
		}

		fmt.Println("🏥 Cluster Health Status:")
		fmt.Println("================================")
		for service, status := range health {
			fmt.Printf("%s: %s\n", strings.Title(service), status)
		}
		return nil
	},
}

func init() {
	// Global k8s flags
	k8sCmd.PersistentFlags().StringVarP(&k8sNamespace, "namespace", "n", "default", "Kubernetes namespace")
	k8sCmd.PersistentFlags().StringVarP(&k8sContext, "context", "c", "", "Kubernetes context")

	// Pods subcommands
	k8sPodsCmd.AddCommand(k8sPodsListCmd)
	k8sPodsLogsCmd.Flags().BoolP("follow", "f", false, "Follow logs")
	k8sPodsCmd.AddCommand(k8sPodsLogsCmd)
	k8sPodsCmd.AddCommand(k8sPodsDescribeCmd)
	k8sPodsCmd.AddCommand(k8sPodsDeleteCmd)

	// Deployments subcommands
	k8sDeploymentsCmd.AddCommand(k8sDeploymentsListCmd)
	k8sDeploymentsCmd.AddCommand(k8sDeploymentsScaleCmd)

	// Services subcommands
	k8sServicesCmd.AddCommand(k8sServicesListCmd)

	// Nodes subcommands
	k8sNodesCmd.AddCommand(k8sNodesListCmd)

	// Context subcommands
	k8sContextCmd.AddCommand(k8sContextListCmd)
	k8sContextCmd.AddCommand(k8sContextCurrentCmd)
	k8sContextCmd.AddCommand(k8sContextSwitchCmd)

	// K8s main subcommands
	k8sCmd.AddCommand(k8sPodsCmd)
	k8sCmd.AddCommand(k8sDeploymentsCmd)
	k8sCmd.AddCommand(k8sServicesCmd)
	k8sCmd.AddCommand(k8sNodesCmd)
	k8sCmd.AddCommand(k8sContextCmd)
	k8sCmd.AddCommand(k8sHealthCmd)
}
