package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	authpkg "github.com/yourusername/ops-tool/pkg/auth"
	azurepkg "github.com/yourusername/ops-tool/pkg/azure"
)

var azureSubscription string
var azureResourceGroup string

var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Microsoft Azure operations",
	Long:  "Manage Azure resources (VMs, Storage, Databases, App Services, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var azureVMCmd = &cobra.Command{
	Use:   "vm",
	Short: "Virtual Machine management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var azureVMListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Azure VMs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := azurepkg.NewClient(azureSubscription, azureResourceGroup)
		output, err := client.ListVMs()
		if err != nil {
			return fmt.Errorf("failed to list VMs: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var azureVMStartCmd = &cobra.Command{
	Use:   "start <vm-name>",
	Short: "Start an Azure VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleOperator); err != nil {
			return err
		}
		client := azurepkg.NewClient(azureSubscription, azureResourceGroup)
		_, err := client.StartVM(args[0])
		if err != nil {
			return fmt.Errorf("failed to start VM: %w", err)
		}
		fmt.Printf("✅ VM %s started\n", args[0])
		return nil
	},
}

var azureVMStopCmd = &cobra.Command{
	Use:   "stop <vm-name>",
	Short: "Stop an Azure VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleOperator); err != nil {
			return err
		}
		client := azurepkg.NewClient(azureSubscription, azureResourceGroup)
		_, err := client.StopVM(args[0])
		if err != nil {
			return fmt.Errorf("failed to stop VM: %w", err)
		}
		fmt.Printf("✅ VM %s stopped\n", args[0])
		return nil
	},
}

var azureStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Storage Account operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var azureStorageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Storage Accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := azurepkg.NewClient(azureSubscription, azureResourceGroup)
		output, err := client.ListStorageAccounts()
		if err != nil {
			return fmt.Errorf("failed to list storage accounts: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var azureDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var azureDatabaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List SQL Databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := azurepkg.NewClient(azureSubscription, azureResourceGroup)
		output, err := client.ListDatabases()
		if err != nil {
			return fmt.Errorf("failed to list databases: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var azureAppServiceCmd = &cobra.Command{
	Use:   "appservice",
	Short: "App Service management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var azureAppServiceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List App Services",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := azurepkg.NewClient(azureSubscription, azureResourceGroup)
		output, err := client.ListAppServices()
		if err != nil {
			return fmt.Errorf("failed to list app services: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var azureResourceGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Resource Group operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var azureResourceGroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Resource Groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := azurepkg.NewClient(azureSubscription, azureResourceGroup)
		output, err := client.ListResourceGroups()
		if err != nil {
			return fmt.Errorf("failed to list resource groups: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var azureNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var azureNetworkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Virtual Networks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := azurepkg.NewClient(azureSubscription, azureResourceGroup)
		output, err := client.ListNetworks()
		if err != nil {
			return fmt.Errorf("failed to list networks: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(azureCmd)

	// VM commands
	azureCmd.AddCommand(azureVMCmd)
	azureVMCmd.AddCommand(azureVMListCmd, azureVMStartCmd, azureVMStopCmd)

	// Storage commands
	azureCmd.AddCommand(azureStorageCmd)
	azureStorageCmd.AddCommand(azureStorageListCmd)

	// Database commands
	azureCmd.AddCommand(azureDatabaseCmd)
	azureDatabaseCmd.AddCommand(azureDatabaseListCmd)

	// App Service commands
	azureCmd.AddCommand(azureAppServiceCmd)
	azureAppServiceCmd.AddCommand(azureAppServiceListCmd)

	// Resource Group commands
	azureCmd.AddCommand(azureResourceGroupCmd)
	azureResourceGroupCmd.AddCommand(azureResourceGroupListCmd)

	// Network commands
	azureCmd.AddCommand(azureNetworkCmd)
	azureNetworkCmd.AddCommand(azureNetworkListCmd)

	// Flags
	azureCmd.PersistentFlags().StringVar(&azureSubscription, "subscription", "", "Azure Subscription ID")
	azureCmd.PersistentFlags().StringVar(&azureResourceGroup, "resource-group", "", "Azure Resource Group")
}
