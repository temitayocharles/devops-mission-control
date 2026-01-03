package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "missionctl",
	Short: "A lightweight, agnostic DevOps CLI tool",
	Long: `missionctl is a unified DevOps command-line interface for managing
Kubernetes, Docker, Cloud infrastructure, Git, and more.

Think of it as k9s for the entire DevOps ecosystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(k8sCmd)
	rootCmd.AddCommand(dockerCmd)
	rootCmd.AddCommand(awsCmd)
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(helpCmd)
	// Register user and login commands if available
	// These are defined in user.go and must be imported for registration
	// The blank import in main.go ensures user.go's init runs

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	// Authorization flags: either supply `--actor <username>` for interactive user
	// or `--token <token>` for API token-based calls. These are used by RBAC checks.
	rootCmd.PersistentFlags().String("actor", "", "actor username performing the action")
	rootCmd.PersistentFlags().String("token", "", "API token to authenticate as a user")
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help for missionctl and all commands",
	Run: func(cmd *cobra.Command, args []string) {
		_ = rootCmd.Help()
	},
}
