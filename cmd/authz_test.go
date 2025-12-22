package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	authpkg "github.com/yourusername/ops-tool/pkg/auth"
)

func TestRequireMinRoleAndResolveActor(t *testing.T) {
	// isolate file-backed stores in a temp dir
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if cerr := os.Chdir(origWd); cerr != nil {
			t.Fatalf("failed to chdir back: %v", cerr)
		}
	})

	// reinitialize stores to use temp dir files
	userStore = authpkg.NewUserStore()
	tokenStore = authpkg.NewTokenStore(filepath.Join(tmp, "tokens.json"))

	// create users
	if err := userStore.AddUser("admin", "pw", authpkg.RoleAdmin); err != nil {
		t.Fatalf("add admin: %v", err)
	}
	if err := userStore.AddUser("bob", "pw", authpkg.RoleViewer); err != nil {
		t.Fatalf("add bob: %v", err)
	}

	// prepare a command and flags
	cmd := &cobra.Command{}
	cmd.Flags().String("actor", "", "actor")
	cmd.Flags().String("token", "", "token")

	// admin should satisfy operator minimum
	if err := cmd.Flags().Set("actor", "admin"); err != nil {
		t.Fatalf("set flag: %v", err)
	}
	if err := requireMinRole(cmd, authpkg.RoleOperator); err != nil {
		t.Fatalf("admin should satisfy operator: %v", err)
	}

	// viewer should NOT satisfy operator
	if err := cmd.Flags().Set("actor", "bob"); err != nil {
		t.Fatalf("set flag: %v", err)
	}
	if err := requireMinRole(cmd, authpkg.RoleOperator); err == nil {
		t.Fatalf("viewer should not satisfy operator")
	}
}

func TestRequireAdminOrSelfAndResolveActorWithToken(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if cerr := os.Chdir(origWd); cerr != nil {
			t.Fatalf("failed to chdir back: %v", cerr)
		}
	})

	userStore = authpkg.NewUserStore()
	tokenStore = authpkg.NewTokenStore(filepath.Join(tmp, "tokens.json"))

	if err := userStore.AddUser("alice", "pw", authpkg.RoleOperator); err != nil {
		t.Fatalf("add alice: %v", err)
	}
	if err := userStore.AddUser("carol", "pw", authpkg.RoleViewer); err != nil {
		t.Fatalf("add carol: %v", err)
	}

	// create a token for alice and ensure resolveActor picks token over actor flag
	tok, err := tokenStore.GenerateToken("alice", "t1", 0)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.Flags().String("actor", "", "actor")
	cmd.Flags().String("token", "", "token")

	if err := cmd.Flags().Set("token", tok.Token); err != nil {
		t.Fatalf("set token flag: %v", err)
	}
	actor, err := resolveActor(cmd)
	if err != nil {
		t.Fatalf("resolveActor error: %v", err)
	}
	if actor != "alice" {
		t.Fatalf("expected actor alice from token, got %q", actor)
	}

	// requireAdminOrSelf: alice acting on herself should pass
	if err := requireAdminOrSelf(cmd, "alice"); err != nil {
		t.Fatalf("alice should be allowed as self: %v", err)
	}

	// carol is viewer; using actor flag should fail for admin-only
	if err := cmd.Flags().Set("token", ""); err != nil {
		t.Fatalf("clear token flag: %v", err)
	}
	if err := cmd.Flags().Set("actor", "carol"); err != nil {
		t.Fatalf("set actor flag: %v", err)
	}
	if err := requireAdmin(cmd); err == nil {
		t.Fatalf("carol is not admin and should be forbidden")
	}
}

func TestRequireMinRole_NoActor(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if cerr := os.Chdir(origWd); cerr != nil {
			t.Fatalf("failed to chdir back: %v", cerr)
		}
	})

	userStore = authpkg.NewUserStore()
	tokenStore = authpkg.NewTokenStore(filepath.Join(tmp, "tokens.json"))

	cmd := &cobra.Command{}
	cmd.Flags().String("actor", "", "actor")
	cmd.Flags().String("token", "", "token")

	if err := requireMinRole(cmd, authpkg.RoleViewer); err == nil {
		t.Fatalf("expected error when no actor or token provided")
	}
}

func TestRequireAdmin_InvalidToken(t *testing.T) {
	origWd, _ := os.Getwd()
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if cerr := os.Chdir(origWd); cerr != nil {
			t.Fatalf("failed to chdir back: %v", cerr)
		}
	}()

	userStore = authpkg.NewUserStore()
	tokenStore = authpkg.NewTokenStore(filepath.Join(tmp, "tokens.json"))

	// set a token that doesn't exist
	cmd := &cobra.Command{}
	cmd.Flags().String("actor", "", "actor")
	cmd.Flags().String("token", "", "token")
	if err := cmd.Flags().Set("token", "deadbeef"); err != nil {
		t.Fatalf("set token: %v", err)
	}
	if _, err := resolveActor(cmd); err == nil {
		t.Fatalf("expected error resolving invalid token")
	}
}
