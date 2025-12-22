package auth

import (
	"os"
	"testing"
	"time"
)

func TestTokenLifecycle(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	us := NewUserStore("")
	if err := us.AddUser("u1", "pw", RoleViewer); err != nil {
		t.Fatal(err)
	}
	ts := NewTokenStore("", "tokens.json")
	tok, err := ts.GenerateToken("u1", "t1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ts.Validate(tok.Token); err != nil {
		t.Fatal(err)
	}
	// rotate
	newTok, err := ts.Rotate(tok.Token, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ts.Validate(tok.Token); err == nil {
		t.Fatalf("old token should be revoked")
	}
	if _, err := ts.Validate(newTok.Token); err != nil {
		t.Fatal(err)
	}
	// revoke
	if err := ts.Revoke(newTok.Token); err != nil {
		t.Fatal(err)
	}
	if _, err := ts.Validate(newTok.Token); err == nil {
		t.Fatalf("revoked token should not validate")
	}
	// expiry
	t2, err := ts.GenerateToken("u1", "t2", 0)
	if err != nil {
		t.Fatal(err)
	}
	// mark expired
	ts.mu.Lock()
	if tt, ok := ts.tokens[t2.Token]; ok {
		past := time.Now().Add(-1 * time.Hour)
		tt.ExpiresAt = &past
	}
	ts.mu.Unlock()
	if _, err := ts.Validate(t2.Token); err == nil {
		t.Fatalf("expired token should not validate")
	}
}
