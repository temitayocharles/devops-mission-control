package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/yourusername/devops-mission-control/pkg/auth"
)

var authCmd = &cobra.Command{
    Use:   "auth",
    Short: "Authentication and user/role management",
}

var authUserCmd = &cobra.Command{
    Use:   "user",
    Short: "User management",
}

var authUserAddCmd = &cobra.Command{
    Use:   "add [username] [role]",
    Short: "Add a user with role",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        username := args[0]
        role := args[1]
        token, err := auth.AddUser(username, role)
        if err != nil {
            return err
        }
        if token == "" {
            fmt.Printf("user %s already exists\n", username)
            return nil
        }
        fmt.Printf("created user %s with role %s\n", username, role)
        fmt.Printf("token: %s\n", token)
        return nil
    },
}

var authUserListCmd = &cobra.Command{
    Use:   "list",
    Short: "List users",
    RunE: func(cmd *cobra.Command, args []string) error {
        users, err := auth.LoadUsers()
        if err != nil {
            return err
        }
        for _, u := range users {
            fmt.Printf("%s\t%s\n", u.Username, u.Role)
        }
        return nil
    },
}

var authUserRemoveCmd = &cobra.Command{
    Use:   "remove [username]",
    Short: "Remove a user",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        username := args[0]
        if err := auth.RemoveUser(username); err != nil {
            return err
        }
        fmt.Printf("removed user %s\n", username)
        return nil
    },
}

var authRoleSetCmd = &cobra.Command{
    Use:   "role set [username] [role]",
    Short: "Set role for a user",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        username := args[0]
        role := args[1]
        if err := auth.SetRole(username, role); err != nil {
            return err
        }
        fmt.Printf("set role %s for %s\n", role, username)
        return nil
    },
}

var authTokenCmd = &cobra.Command{
    Use:   "token",
    Short: "Token management",
}

var authTokenListCmd = &cobra.Command{
    Use:   "list",
    Short: "List tokens (username:token)",
    RunE: func(cmd *cobra.Command, args []string) error {
        users, err := auth.LoadUsers()
        if err != nil {
            return err
        }
        for _, u := range users {
            fmt.Printf("%s\t%s\n", u.Username, u.Token)
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(authCmd)
    authCmd.AddCommand(authUserCmd)
    authUserCmd.AddCommand(authUserAddCmd)
    authUserCmd.AddCommand(authUserListCmd)
    authUserCmd.AddCommand(authUserRemoveCmd)
    authCmd.AddCommand(authRoleSetCmd)
    authCmd.AddCommand(authTokenCmd)
    authTokenCmd.AddCommand(authTokenListCmd)

    // ensure users.json exists
    if _, err := os.Stat("users.json"); os.IsNotExist(err) {
        _ = auth.SaveUsers([]auth.User{})
    }
}
