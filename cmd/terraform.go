package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	authpkg "github.com/yourusername/devops-mission-control/pkg/auth"
	terraformpkg "github.com/yourusername/devops-mission-control/pkg/terraform"
)

var terraformWorkdir string

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Terraform infrastructure management",
	Long:  "Manage infrastructure using Terraform (plan, apply, destroy, state, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var terraformVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get Terraform version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.Version()
		if err != nil {
			return fmt.Errorf("failed to get version: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Terraform configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Terraform working directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleOperator); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.Init()
		if err != nil {
			return fmt.Errorf("failed to init: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Create Terraform execution plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleOperator); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.Plan("")
		if err != nil {
			return fmt.Errorf("failed to plan: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Terraform changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleAdmin); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.Apply("")
		if err != nil {
			return fmt.Errorf("failed to apply: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy Terraform-managed infrastructure",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.Destroy("")
		if err != nil {
			return fmt.Errorf("failed to destroy: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformFormatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format Terraform configuration files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.Format()
		if err != nil {
			return fmt.Errorf("failed to format: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Terraform state management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var terraformStateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Terraform state resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.StateList()
		if err != nil {
			return fmt.Errorf("failed to list state: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformStateShowCmd = &cobra.Command{
	Use:   "show <resource>",
	Short: "Show Terraform state resource",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.StateShow(args[0])
		if err != nil {
			return fmt.Errorf("failed to show state: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformOutputCmd = &cobra.Command{
	Use:   "output",
	Short: "Show Terraform outputs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.OutputAll()
		if err != nil {
			return fmt.Errorf("failed to get output: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

var terraformWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Terraform workspace management",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

var terraformWorkspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Terraform workspaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		client := terraformpkg.NewClient(terraformWorkdir)
		output, err := client.WorkspaceList()
		if err != nil {
			return fmt.Errorf("failed to list workspaces: %w", err)
		}
		fmt.Println(output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)
	terraformCmd.AddCommand(
		terraformVersionCmd,
		terraformValidateCmd,
		terraformInitCmd,
		terraformPlanCmd,
		terraformApplyCmd,
		terraformDestroyCmd,
		terraformFormatCmd,
		terraformOutputCmd,
	)

	// State commands
	terraformCmd.AddCommand(terraformStateCmd)
	terraformStateCmd.AddCommand(terraformStateListCmd, terraformStateShowCmd)

	// Workspace commands
	terraformCmd.AddCommand(terraformWorkspaceCmd)
	terraformWorkspaceCmd.AddCommand(terraformWorkspaceListCmd)

	// Flags
	terraformCmd.PersistentFlags().StringVarP(&terraformWorkdir, "workdir", "d", ".", "Terraform working directory")
}
