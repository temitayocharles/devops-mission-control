package dashboard

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	authpkg "github.com/yourusername/ops-tool/pkg/auth"
)

// Test that the dashboard accepts a token created after the server started.
func TestDashboardPicksUpRuntimeToken(t *testing.T) {
	// use temp dir so token store writes to a temp file
	td, err := os.MkdirTemp("", "dash-integ")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(td); err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}()

	// chdir to temp dir so TokenStore uses that location
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(td); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if cerr := os.Chdir(oldwd); cerr != nil {
			t.Fatalf("failed to chdir back: %v", cerr)
		}
	}()

	// create user store and add a user (NewUserStore has no args now)
	us := authpkg.NewUserStore()
	if err := us.AddUser("alice", "password", authpkg.RoleViewer); err != nil {
		t.Fatal(err)
	}
	httpUserStore = us

	// start a token store (this will create tokens.json file path)
	ts := authpkg.NewTokenStore("tokens.json")
	httpTokenStore = ts

	// start HTTP handler using authMiddleware (viewer for reads)
	handler := authMiddleware(authpkg.RoleViewer, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// first, request should fail as there is no token
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected unauthorized before token exists, got 200")
	}

	// create a token using a separate TokenStore instance that writes to the same file
	ts2 := authpkg.NewTokenStore("tokens.json")
	tok, err := ts2.GenerateToken("alice", "runtime-token", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// small pause to allow file write to flush
	time.Sleep(100 * time.Millisecond)

	// request again with Authorization header using the newly created token
	req2 := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	req2.Header.Set("Authorization", "Bearer "+tok.Token)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200 after token created, got %d", rr2.Code)
	}

	// sanity: ensure tokens.json exists
	if _, err := os.Stat(filepath.Join(td, "tokens.json")); err != nil {
		t.Fatalf("tokens.json not found: %v", err)
	}
}
