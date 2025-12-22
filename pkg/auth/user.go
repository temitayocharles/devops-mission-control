package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/yourusername/ops-tool/pkg/audit"
)

// Role defines user access level
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

// User represents a system user
type User struct {
	Username     string    `json:"username"`
	TenantID     string    `json:"tenant_id,omitempty"`
	PasswordHash string    `json:"password_hash"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	Active       bool      `json:"active"`
}

// UserStore manages users and roles
type UserStore struct {
	mu    sync.RWMutex
	users map[string]*User
	file  string
	tenant string
}

// NewUserStore creates a new user store
func NewUserStore(tenant string) *UserStore {
	file := "users.json"
	if tenant != "" {
		file = fmt.Sprintf("users.%s.json", tenant)
	}
	store := &UserStore{
		users:  make(map[string]*User),
		file:   file,
		tenant: tenant,
	}
	store.load()
	return store
}

// AddUser adds a new user
func (s *UserStore) AddUser(username, password string, role Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[username]; exists {
		return errors.New("user already exists")
	}
	hash := hashPassword(password)
	s.users[username] = &User{
		Username:     username,
		TenantID:     s.tenant,
		PasswordHash: hash,
		Role:         role,
		CreatedAt:    time.Now(),
		Active:       true,
	}
	if err := s.save(); err != nil {
		return err
	}
	if err := audit.Record(s.tenant, "user.create", username, username, map[string]any{"role": role}); err != nil {
		fmt.Fprintf(os.Stderr, "audit record failed: %v\n", err)
	}
	return nil
}

// Authenticate checks user credentials
func (s *UserStore) Authenticate(username, password string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[username]
	if !ok || !user.Active {
		return nil, errors.New("invalid user or inactive")
	}
	if user.PasswordHash != hashPassword(password) {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

// GetUser returns a user by username
func (s *UserStore) GetUser(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// ListUsers returns all users
func (s *UserStore) ListUsers() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

// SetUserActive sets a user's active status
func (s *UserStore) SetUserActive(username string, active bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[username]
	if !ok {
		return errors.New("user not found")
	}
	user.Active = active
	if err := s.save(); err != nil {
		return err
	}
	if err := audit.Record(s.tenant, "user.set_active", username, username, map[string]any{"active": active}); err != nil {
		fmt.Fprintf(os.Stderr, "audit record failed: %v\n", err)
	}
	return nil
}

// SetUserRole updates a user's role
func (s *UserStore) SetUserRole(username string, role Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[username]
	if !ok {
		return errors.New("user not found")
	}
	user.Role = role
	if err := s.save(); err != nil {
		return err
	}
	if err := audit.Record(s.tenant, "user.set_role", username, username, map[string]any{"role": role}); err != nil {
		fmt.Fprintf(os.Stderr, "audit record failed: %v\n", err)
	}
	return nil
}

// DeleteUser removes a user from the store
func (s *UserStore) DeleteUser(username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[username]; !ok {
		return errors.New("user not found")
	}
	delete(s.users, username)
	if err := s.save(); err != nil {
		return err
	}
	if err := audit.Record(s.tenant, "user.delete", username, username, nil); err != nil {
		fmt.Fprintf(os.Stderr, "audit record failed: %v\n", err)
	}
	return nil
}

// hashPassword hashes a password
func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

// save persists users to disk
func (s *UserStore) save() error {
	// Note: callers (AddUser, SetUserActive, DeleteUser) hold s.mu.Lock()
	// so `save` must not attempt to acquire the mutex again (would deadlock).
	f, err := os.Create(s.file)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close users file: %v\n", cerr)
		}
	}()
	enc := json.NewEncoder(f)
	return enc.Encode(s.users)
}

// load loads users from disk
func (s *UserStore) load() {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.Open(s.file)
	if err != nil {
		return
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close users file: %v\n", cerr)
		}
	}()
	dec := json.NewDecoder(f)
	_ = dec.Decode(&s.users)
}
