package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	gcppkg "github.com/yourusername/ops-tool/pkg/gcp"
)

var gcpProject string
var gcpRegion string

var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Google Cloud Platform operations",
	Long:  "Manage Google Cloud resources (Compute Engine, Cloud Storage, Cloud SQL, IAM, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gcpComputeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Compute Engine instance management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gcpComputeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Compute Engine instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gcppkg.NewClient(gcpProject, gcpRegion)
		output, err := client.ListInstances()
		if err != nil {
			return fmt.Errorf("failed to list instances: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var gcpComputeStartCmd = &cobra.Command{
	Use:   "start <instance>",
	Short: "Start a Compute Engine instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gcppkg.NewClient(gcpProject, gcpRegion)
		_, err := client.StartInstance(args[0])
		if err != nil {
			return fmt.Errorf("failed to start instance: %w", err)
		}
		fmt.Printf("✅ Instance %s started\n", args[0])
		return nil
	},
}

var gcpComputeStopCmd = &cobra.Command{
	Use:   "stop <instance>",
	Short: "Stop a Compute Engine instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gcppkg.NewClient(gcpProject, gcpRegion)
		_, err := client.StopInstance(args[0])
		if err != nil {
			return fmt.Errorf("failed to stop instance: %w", err)
		}
		fmt.Printf("✅ Instance %s stopped\n", args[0])
		return nil
	},
}

var gcpStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Cloud Storage bucket operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gcpStorageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Cloud Storage buckets",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gcppkg.NewClient(gcpProject, gcpRegion)
		output, err := client.ListBuckets()
		if err != nil {
			return fmt.Errorf("failed to list buckets: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var gcpSQLCmd = &cobra.Command{
	Use:   "sql",
	Short: "Cloud SQL database management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gcpSQLListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Cloud SQL instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gcppkg.NewClient(gcpProject, gcpRegion)
		output, err := client.ListSQLInstances()
		if err != nil {
			return fmt.Errorf("failed to list SQL instances: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var gcpRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Cloud Run service management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gcpRunListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Cloud Run services",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gcppkg.NewClient(gcpProject, gcpRegion)
		output, err := client.ListCloudRunServices()
		if err != nil {
			return fmt.Errorf("failed to list Cloud Run services: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var gcpIAMCmd = &cobra.Command{
	Use:   "iam",
	Short: "IAM and access control",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var gcpIAMAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List service accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gcppkg.NewClient(gcpProject, gcpRegion)
		output, err := client.ListServiceAccounts()
		if err != nil {
			return fmt.Errorf("failed to list service accounts: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var gcpInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display GCP project information",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := gcppkg.NewClient(gcpProject, gcpRegion)
		output, err := client.GetProjectInfo()
		if err != nil {
			return fmt.Errorf("failed to get project info: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gcpCmd)
	
	// Compute commands
	gcpCmd.AddCommand(gcpComputeCmd)
	gcpComputeCmd.AddCommand(gcpComputeListCmd, gcpComputeStartCmd, gcpComputeStopCmd)
	
	// Storage commands
	gcpCmd.AddCommand(gcpStorageCmd)
	gcpStorageCmd.AddCommand(gcpStorageListCmd)
	
	// SQL commands
	gcpCmd.AddCommand(gcpSQLCmd)
	gcpSQLCmd.AddCommand(gcpSQLListCmd)
	
	// Cloud Run commands
	gcpCmd.AddCommand(gcpRunCmd)
	gcpRunCmd.AddCommand(gcpRunListCmd)
	
	// IAM commands
	gcpCmd.AddCommand(gcpIAMCmd)
	gcpIAMCmd.AddCommand(gcpIAMAccountsCmd)
	
	// Project info
	gcpCmd.AddCommand(gcpInfoCmd)
	
	// Flags
	gcpCmd.PersistentFlags().StringVar(&gcpProject, "project", "", "GCP project ID")
	gcpCmd.PersistentFlags().StringVar(&gcpRegion, "region", "us-central1", "GCP region")
}
