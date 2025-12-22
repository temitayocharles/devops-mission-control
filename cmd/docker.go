package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/ops-tool/pkg/docker"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker container operations",
	Long:  "Manage Docker containers, images, and compose",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var dockerContainersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Manage containers",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var dockerContainersListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List running containers",
	Aliases: []string{"ls", "ps"},
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		client := docker.NewClient()
		output, err := client.ListContainers(!all)
		if err != nil {
			return fmt.Errorf("failed to list containers: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var dockerContainersStopCmd = &cobra.Command{
	Use:   "stop <container-id|name>",
	Short: "Stop a running container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := docker.NewClient()
		output, err := client.StopContainer(args[0])
		if err != nil {
			return fmt.Errorf("failed to stop container: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var dockerContainersRemoveCmd = &cobra.Command{
	Use:   "remove <container-id|name>",
	Short: "Remove a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		client := docker.NewClient()
		output, err := client.RemoveContainer(args[0], force)
		if err != nil {
			return fmt.Errorf("failed to remove container: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var dockerContainersLogsCmd = &cobra.Command{
	Use:   "logs <container-id|name>",
	Short: "View container logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		follow, _ := cmd.Flags().GetBool("follow")
		client := docker.NewClient()
		return client.GetContainerLogs(args[0], follow)
	},
}

var dockerContainersStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show container resource usage",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := docker.NewClient()
		output, err := client.GetContainerStats()
		if err != nil {
			return fmt.Errorf("failed to get stats: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var dockerImagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Manage Docker images",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var dockerImagesListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Docker images",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client := docker.NewClient()
		output, err := client.ListImages()
		if err != nil {
			return fmt.Errorf("failed to list images: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var dockerImagesPullCmd = &cobra.Command{
	Use:   "pull <image>",
	Short: "Pull an image from registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := docker.NewClient()
		return client.PullImage(args[0])
	},
}

var dockerImagesRemoveCmd = &cobra.Command{
	Use:   "remove <image-id>",
	Short: "Remove an image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		client := docker.NewClient()
		output, err := client.RemoveImage(args[0], force)
		if err != nil {
			return fmt.Errorf("failed to remove image: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var dockerComposeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Docker Compose operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var dockerComposeUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start Docker Compose services",
	RunE: func(cmd *cobra.Command, args []string) error {
		detach, _ := cmd.Flags().GetBool("detach")
		client := docker.NewClient()
		return client.ComposeUp(detach)
	},
}

var dockerComposeDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop Docker Compose services",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := docker.NewClient()
		return client.ComposeDown()
	},
}

var dockerComposeLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View Docker Compose logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		follow, _ := cmd.Flags().GetBool("follow")
		client := docker.NewClient()
		return client.ComposeLogs(follow)
	},
}

var dockerComposeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Docker Compose service status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := docker.NewClient()
		output, err := client.ComposeStatus()
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var dockerSystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Docker system operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var dockerSystemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show Docker system information",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := docker.NewClient()
		output, err := client.GetSystemInfo()
		if err != nil {
			return fmt.Errorf("failed to get system info: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

var dockerSystemPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove unused Docker resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		client := docker.NewClient()
		output, err := client.SystemPrune(all)
		if err != nil {
			return fmt.Errorf("failed to prune: %w", err)
		}

		fmt.Println(output)
		return nil
	},
}

func init() {
	// Containers subcommands
	dockerContainersListCmd.Flags().BoolP("all", "a", false, "Show all containers")
	dockerContainersCmd.AddCommand(dockerContainersListCmd)

	dockerContainersRemoveCmd.Flags().BoolP("force", "f", false, "Force remove")
	dockerContainersCmd.AddCommand(dockerContainersRemoveCmd)

	dockerContainersStopCmd.Flags().BoolP("force", "f", false, "Force stop")
	dockerContainersCmd.AddCommand(dockerContainersStopCmd)

	dockerContainersLogsCmd.Flags().BoolP("follow", "f", false, "Follow logs")
	dockerContainersCmd.AddCommand(dockerContainersLogsCmd)

	dockerContainersCmd.AddCommand(dockerContainersStatsCmd)

	// Images subcommands
	dockerImagesCmd.AddCommand(dockerImagesListCmd)
	dockerImagesCmd.AddCommand(dockerImagesPullCmd)

	dockerImagesRemoveCmd.Flags().BoolP("force", "f", false, "Force remove")
	dockerImagesCmd.AddCommand(dockerImagesRemoveCmd)

	// Compose subcommands
	dockerComposeUpCmd.Flags().BoolP("detach", "d", false, "Detached mode")
	dockerComposeCmd.AddCommand(dockerComposeUpCmd)
	dockerComposeCmd.AddCommand(dockerComposeDownCmd)

	dockerComposeLogsCmd.Flags().BoolP("follow", "f", false, "Follow logs")
	dockerComposeCmd.AddCommand(dockerComposeLogsCmd)

	dockerComposeCmd.AddCommand(dockerComposeStatusCmd)

	// System subcommands
	dockerSystemCmd.AddCommand(dockerSystemInfoCmd)

	dockerSystemPruneCmd.Flags().BoolP("all", "a", false, "Remove all unused resources")
	dockerSystemCmd.AddCommand(dockerSystemPruneCmd)

	// Docker main subcommands
	dockerCmd.AddCommand(dockerContainersCmd)
	dockerCmd.AddCommand(dockerImagesCmd)
	dockerCmd.AddCommand(dockerComposeCmd)
	dockerCmd.AddCommand(dockerSystemCmd)
}
