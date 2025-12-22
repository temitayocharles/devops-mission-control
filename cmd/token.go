package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	authpkg "github.com/yourusername/ops-tool/pkg/auth"
)

var tokenStore = authpkg.NewTokenStore("tokens.json")

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "API token management",
}

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an API token for a user",
	Args:  cobra.ExactArgs(2), // username, name
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		name := args[1]
		ttlHours, _ := cmd.Flags().GetInt("ttl")
		duration := time.Duration(ttlHours) * time.Hour
		// verify user exists
		if _, err := userStore.GetUser(username); err != nil {
			return err
		}
		// allow creating a token for yourself or require admin
		if err := requireAdminOrSelf(cmd, username); err != nil {
			return err
		}
		tok, err := tokenStore.GenerateToken(username, name, duration)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Token created: %s\n", tok.Token)
		return nil
	},
}

var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		// If actor is admin, show all tokens; otherwise show only tokens belonging to actor
		actor, _ := resolveActor(cmd)
		toks := tokenStore.ListTokens()
		filtered := make([]*authpkg.Token, 0, len(toks))
		if actor == "" {
			// no actor supplied -> require admin
			if err := requireAdmin(cmd); err != nil {
				return fmt.Errorf("forbidden: provide --actor or --token")
			}
			filtered = toks
		} else {
			// check role; if admin see all, else filter
			u, err := userStore.GetUser(actor)
			if err == nil && u.Role == authpkg.RoleAdmin {
				filtered = toks
			} else {
				for _, t := range toks {
					if t.User == actor {
						filtered = append(filtered, t)
					}
				}
			}
		}
		if len(filtered) == 0 {
			fmt.Println("No tokens found")
			return nil
		}
		fmt.Println("Tokens:")
		for _, t := range filtered {
			exp := "never"
			if t.ExpiresAt != nil {
				exp = t.ExpiresAt.String()
			}
			fmt.Printf("  %s (user: %s, name: %s, revoked: %v, expires: %s)\n", t.Token, t.User, t.Name, t.Revoked, exp)
		}
		return nil
	},
}

var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke an API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t := args[0]
		// only admin or token owner may revoke
		tok, err := tokenStore.Validate(t)
		if err != nil {
			return err
		}
		if err := requireAdminOrSelf(cmd, tok.User); err != nil {
			return err
		}
		if err := tokenStore.Revoke(t); err != nil {
			return err
		}
		fmt.Println("✅ Token revoked")
		return nil
	},
}

var tokenRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate an API token (issue new, revoke old)",
	Args:  cobra.ExactArgs(1), // old token
	RunE: func(cmd *cobra.Command, args []string) error {
		old := args[0]
		tok, err := tokenStore.Validate(old)
		if err != nil {
			return err
		}
		// only owner or admin may rotate
		if err := requireAdminOrSelf(cmd, tok.User); err != nil {
			return err
		}
		ttlHours, _ := cmd.Flags().GetInt("ttl")
		duration := time.Duration(ttlHours) * time.Hour
		newTok, err := tokenStore.Rotate(old, duration)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Token rotated: %s (replaced %s)\n", newTok.Token, old)
		return nil
	},
}

var tokenValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireMinRole(cmd, authpkg.RoleViewer); err != nil {
			return err
		}
		t := args[0]
		tok, err := tokenStore.Validate(t)
		if err != nil {
			return err
		}
		exp := "never"
		if tok.ExpiresAt != nil {
			exp = tok.ExpiresAt.String()
		}
		fmt.Printf("Token valid: user=%s name=%s expires=%s revoked=%v\n", tok.User, tok.Name, exp, tok.Revoked)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.AddCommand(tokenCreateCmd, tokenListCmd, tokenRevokeCmd, tokenValidateCmd)
	tokenCreateCmd.Flags().Int("ttl", 0, "TTL in hours for the token (0 = no expiry)")
	tokenRotateCmd.Flags().Int("ttl", 0, "TTL in hours for the new token (0 = no expiry)")
	tokenCmd.AddCommand(tokenRotateCmd)
}
