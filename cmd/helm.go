package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	authpkg "github.com/yourusername/ops-tool/pkg/auth"
	helmpkg "github.com/yourusername/ops-tool/pkg/helm"
)

var helmNamespace string

var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm chart management",
	Long:  "Manage Kubernetes Helm charts and releases",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var helmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Helm releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := helmpkg.NewClient(helmNamespace)
		output, err := client.ListReleases()
		if err != nil {
			return fmt.Errorf("failed to list releases: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var helmRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Helm repository operations",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var helmRepoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Helm repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := helmpkg.NewClient(helmNamespace)
		output, err := client.RepoList()
		if err != nil {
			return fmt.Errorf("failed to list repos: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var helmRepoUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Helm repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := helmpkg.NewClient(helmNamespace)
		output, err := client.RepoUpdate()
		if err != nil {
			return fmt.Errorf("failed to update repos: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var helmSearchCmd = &cobra.Command{
	Use:   "search <chart-name>",
	Short: "Search for Helm charts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := helmpkg.NewClient(helmNamespace)
		output, err := client.SearchChart(args[0])
		if err != nil {
			return fmt.Errorf("failed to search charts: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var helmStatusCmd = &cobra.Command{
	Use:   "status <release-name>",
	Short: "Get Helm release status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := helmpkg.NewClient(helmNamespace)
		output, err := client.GetReleaseStatus(args[0])
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var helmHistoryCmd = &cobra.Command{
	Use:   "history <release-name>",
	Short: "Get Helm release history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := helmpkg.NewClient(helmNamespace)
		output, err := client.GetReleaseHistory(args[0])
		if err != nil {
			return fmt.Errorf("failed to get history: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var helmValuesCmd = &cobra.Command{
	Use:   "values <release-name>",
	Short: "Get Helm release values",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := helmpkg.NewClient(helmNamespace)
		output, err := client.GetReleaseValues(args[0])
		if err != nil {
			return fmt.Errorf("failed to get values: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var helmUninstallCmd = &cobra.Command{
	Use:   "uninstall <release-name>",
	Short: "Uninstall a Helm release",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleOperator); err != nil {
			return err
		}
		client := helmpkg.NewClient(helmNamespace)
		_, err := client.UninstallRelease(args[0])
		if err != nil {
			return fmt.Errorf("failed to uninstall release: %w", err)
		}
		fmt.Printf("✅ Release %s uninstalled\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(helmCmd)
	helmCmd.AddCommand(helmListCmd, helmSearchCmd, helmStatusCmd, helmHistoryCmd, helmValuesCmd, helmUninstallCmd)

	// Repo commands
	helmCmd.AddCommand(helmRepoCmd)
	helmRepoCmd.AddCommand(helmRepoListCmd, helmRepoUpdateCmd)

	// Flags
	helmCmd.PersistentFlags().StringVarP(&helmNamespace, "namespace", "n", "default", "Kubernetes namespace")
}
