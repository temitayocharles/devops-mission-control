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
	Run: func(cmd *cobra.Command, args []string) {
		toks := tokenStore.ListTokens()
		if len(toks) == 0 {
			fmt.Println("No tokens found")
			return
		}
		fmt.Println("Tokens:")
		for _, t := range toks {
			exp := "never"
			if t.ExpiresAt != nil {
				exp = t.ExpiresAt.String()
			}
			fmt.Printf("  %s (user: %s, name: %s, revoked: %v, expires: %s)\n", t.Token, t.User, t.Name, t.Revoked, exp)
		}
	},
}

var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke an API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t := args[0]
		if err := tokenStore.Revoke(t); err != nil {
			return err
		}
		fmt.Println("✅ Token revoked")
		return nil
	},
}

var tokenValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
}
