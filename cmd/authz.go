package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	authpkg "github.com/yourusername/devops-mission-control/pkg/auth"
)

// resolveActor determines the actor username from flags: token takes precedence
// over the explicit --actor flag. Returns empty string if none provided.
func resolveActor(cmd *cobra.Command) (string, error) {
	tokenStr, _ := cmd.Flags().GetString("token")
	if tokenStr != "" {
		tok, err := tokenStore.Validate(tokenStr)
		if err != nil {
			return "", fmt.Errorf("invalid token: %w", err)
		}
		return tok.User, nil
	}
	actor, _ := cmd.Flags().GetString("actor")
	if actor != "" {
		return actor, nil
	}
	return "", nil
}

// requireAdmin ensures the resolved actor is an admin user.
func requireAdmin(cmd *cobra.Command) error {
	actor, err := resolveActor(cmd)
	if err != nil {
		return err
	}
	if actor == "" {
		return errors.New("no actor provided; use --actor or --token")
	}
	u, err := userStore.GetUser(actor)
	if err != nil {
		return fmt.Errorf("actor lookup failed: %w", err)
	}
	if u.Role != authpkg.RoleAdmin {
		return errors.New("forbidden: admin role required")
	}
	return nil
}

// requireAdminOrSelf allows action if actor is admin or actor equals targetUsername
func requireAdminOrSelf(cmd *cobra.Command, targetUsername string) error {
	actor, err := resolveActor(cmd)
	if err != nil {
		return err
	}
	if actor == "" {
		return errors.New("no actor provided; use --actor or --token")
	}
	if actor == targetUsername {
		return nil
	}
	return requireAdmin(cmd)
}

// role level mapping used by CLI RBAC checks
var cliRoleLevel = map[authpkg.Role]int{
	authpkg.RoleViewer:   10,
	authpkg.RoleOperator: 20,
	authpkg.RoleAdmin:    30,
}

// requireMinRole enforces that the resolved actor has at least the provided role.
func requireMinRole(cmd *cobra.Command, min authpkg.Role) error {
	actor, err := resolveActor(cmd)
	if err != nil {
		return err
	}
	if actor == "" {
		return errors.New("no actor provided; use --actor or --token")
	}
	u, err := userStore.GetUser(actor)
	if err != nil {
		return fmt.Errorf("actor lookup failed: %w", err)
	}
	if cliRoleLevel[u.Role] < cliRoleLevel[min] {
		return errors.New("forbidden: insufficient role")
	}
	return nil
}
