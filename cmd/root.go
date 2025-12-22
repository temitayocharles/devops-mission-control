package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ops-tool",
	Short: "A lightweight, agnostic DevOps CLI tool",
	Long: `ops-tool is a unified DevOps command-line interface for managing
Kubernetes, Docker, Cloud infrastructure, Git, and more.

Think of it as k9s for the entire DevOps ecosystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
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
	
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help for ops-tool and all commands",
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.Help()
	},
}
