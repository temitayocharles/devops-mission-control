package dashboard

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	authpkg "github.com/yourusername/ops-tool/pkg/auth"
)

func TestAuthMiddleware_AllowsViewerWithToken(t *testing.T) {
	dir := t.TempDir()
	// switch cwd so stores use temp files
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if cerr := os.Chdir(cwd); cerr != nil {
			t.Fatalf("failed to chdir back: %v", cerr)
		}
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	us := authpkg.NewUserStore()
	if err := us.AddUser("alice", "pw", authpkg.RoleViewer); err != nil {
		t.Fatal(err)
	}
	ts := authpkg.NewTokenStore("tokens.json")
	tok, err := ts.GenerateToken("alice", "test", 0)
	if err != nil {
		t.Fatal(err)
	}

	// point package stores at our temp stores
	httpUserStore = us
	httpTokenStore = ts

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/metrics", nil)
	req.Header.Set("Authorization", "Bearer "+tok.Token)

	handler := authMiddleware(authpkg.RoleViewer, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		if _, werr := w.Write([]byte("ok")); werr != nil {
			t.Fatalf("write failed: %v", werr)
		}
	})

	handler.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("expected 200, got %d: body=%s", rr.Code, rr.Body.String())
	}
}

func TestAuthMiddleware_RejectsLowerRole(t *testing.T) {
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if cerr := os.Chdir(cwd); cerr != nil {
			t.Fatalf("failed to chdir back: %v", cerr)
		}
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	us := authpkg.NewUserStore()
	if err := us.AddUser("bob", "pw", authpkg.RoleViewer); err != nil {
		t.Fatal(err)
	}
	httpUserStore = us

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/stats", nil)
	req.Header.Set("X-Actor", "bob")

	handler := authMiddleware(authpkg.RoleAdmin, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}
