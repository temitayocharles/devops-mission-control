package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	authpkg "github.com/yourusername/ops-tool/pkg/auth"
)

var userStore = authpkg.NewUserStore()

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
}

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		username, password, role := args[0], args[1], args[2]
		// Only admin may create arbitrary users; allow creating self if actor==username
		if err := requireAdminOrSelf(cmd, username); err != nil {
			return err
		}
		err := userStore.AddUser(username, password, authpkg.Role(role))
		if err != nil {
			return err
		}
		fmt.Printf("✅ User '%s' created with role '%s'\n", username, role)
		return nil
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		users := userStore.ListUsers()
		if len(users) == 0 {
			fmt.Println("No users found")
			return nil
		}
		fmt.Println("Users:")
		for _, u := range users {
			fmt.Printf("  %s (role: %s, active: %v)\n", u.Username, u.Role, u.Active)
		}
		return nil
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		if err := requireAdmin(cmd); err != nil {
			return err
		}
		err := userStore.DeleteUser(username)
		if err != nil {
			return err
		}
		fmt.Printf("✅ User '%s' deleted\n", username)
		return nil
	},
}

var userSetRoleCmd = &cobra.Command{
	Use:   "set-role",
	Short: "Set a user's role",
	Args:  cobra.ExactArgs(2), // username, role
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		role := args[1]
		if err := requireAdmin(cmd); err != nil {
			return err
		}
		if err := userStore.SetUserRole(username, authpkg.Role(role)); err != nil {
			return err
		}
		fmt.Printf("✅ User '%s' role set to '%s'\n", username, role)
		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate as a user",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		username, password := args[0], args[1]
		user, err := userStore.Authenticate(username, password)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Logged in as '%s' (role: %s)\n", user.Username, user.Role)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userCreateCmd, userListCmd, userDeleteCmd, userSetRoleCmd)
	rootCmd.AddCommand(loginCmd)
}
