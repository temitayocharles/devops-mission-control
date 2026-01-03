package auth

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "errors"
	type Role string
	type User struct {
	type UserStore struct {
	package auth

	import (
		"crypto/rand"
		"encoding/hex"
		"encoding/json"
		"errors"
		"log"
		"os"
		"sync"
		"time"
	)

	// Role defines simple RBAC roles
	type Role string

	const (
		RoleViewer   Role = "viewer"
		RoleOperator Role = "operator"
		RoleAdmin    Role = "admin"
	)

	// User represents a local user for RBAC/auth
	type User struct {
		Username  string    `json:"username"`
		Password  string    `json:"password,omitempty"`
		Role      Role      `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	// UserStore manages users persisted to a JSON file.
	type UserStore struct {
		mu        sync.RWMutex
		users     map[string]*User
		usersFile string
	}

	// NewUserStore creates a user store backed by users.json in dir (or cwd if empty)
	func NewUserStore(dir string) *UserStore {
		us := &UserStore{users: make(map[string]*User)}
		if dir == "" {
			us.usersFile = "users.json"
		} else {
			us.usersFile = dir + "/users.json"
		}
		_ = us.load()
		return us
	}

	func (us *UserStore) load() error {
		us.mu.Lock()
		defer us.mu.Unlock()
		f, err := os.Open(us.usersFile)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		defer f.Close()
		var list []*User
		dec := json.NewDecoder(f)
		if err := dec.Decode(&list); err != nil {
			return err
		}
		for _, u := range list {
			us.users[u.Username] = u
		}
		return nil
	}

	func (us *UserStore) save() error {
		us.mu.RLock()
		list := make([]*User, 0, len(us.users))
		for _, u := range us.users {
			list = append(list, u)
		}
		us.mu.RUnlock()

		f, err := os.Create(us.usersFile)
		if err != nil {
			return err
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(list)
	}

	// AddUser adds a user with password and role
	func (us *UserStore) AddUser(username, password string, role Role) error {
		us.mu.Lock()
		defer us.mu.Unlock()
		if _, ok := us.users[username]; ok {
			return nil
		}
		u := &User{Username: username, Password: password, Role: role, CreatedAt: time.Now()}
		us.users[username] = u
		return us.save()
	}

	// GetUser returns a user or error
	func (us *UserStore) GetUser(username string) (*User, error) {
		us.mu.RLock()
		defer us.mu.RUnlock()
		if u, ok := us.users[username]; ok {
			return u, nil
		}
		return nil, os.ErrNotExist
	}

	// ListUsers returns all users
	func (us *UserStore) ListUsers() []*User {
		us.mu.RLock()
		defer us.mu.RUnlock()
		out := make([]*User, 0, len(us.users))
		for _, u := range us.users {
			out = append(out, u)
		}
		return out
	}

	// DeleteUser removes a user
	func (us *UserStore) DeleteUser(username string) error {
		us.mu.Lock()
		defer us.mu.Unlock()
		delete(us.users, username)
		return us.save()
	}

	// SetUserRole updates user's role.
	func (us *UserStore) SetUserRole(username string, role Role) error {
		us.mu.Lock()
		defer us.mu.Unlock()
		u, ok := us.users[username]
		if !ok {
			return errors.New("user not found")
		}
		u.Role = role
		return us.save()
	}

	// Authenticate validates username/password and returns the user.
	func (us *UserStore) Authenticate(username, password string) (*User, error) {
		us.mu.RLock()
		defer us.mu.RUnlock()
		u, ok := us.users[username]
		if !ok {
			return nil, errors.New("invalid credentials")
		}
		if u.Password != password {
			return nil, errors.New("invalid credentials")
		}
		return u, nil
	}

	// Token represents an issued API token
	type Token struct {
		Token     string     `json:"token"`
		User      string     `json:"user"`
		Name      string     `json:"name"`
		ExpiresAt *time.Time `json:"expires_at,omitempty"`
		Revoked   bool       `json:"revoked"`
		CreatedAt time.Time  `json:"created_at"`
	}

	// TokenStore manages tokens persisted to tokens.json
	type TokenStore struct {
		mu         sync.RWMutex
		tokens     map[string]*Token
		tokensFile string
		stopChan   chan struct{}
	}

	// NewTokenStore creates a token store backed by tokens.json in dir (or cwd if empty)
	func NewTokenStore(dir, tokensFile string) *TokenStore {
		ts := &TokenStore{tokens: make(map[string]*Token), stopChan: make(chan struct{})}
		if tokensFile == "" {
			ts.tokensFile = "tokens.json"
		} else {
			ts.tokensFile = tokensFile
		}
		_ = ts.load()
		return ts
	}

	func (ts *TokenStore) load() error {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		f, err := os.Open(ts.tokensFile)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		defer f.Close()
		var list []*Token
		dec := json.NewDecoder(f)
		if err := dec.Decode(&list); err != nil {
			return err
		}
		for _, t := range list {
			ts.tokens[t.Token] = t
		}
		return nil
	}

	func (ts *TokenStore) persistList() error {
		ts.mu.RLock()
		copyList := make([]*Token, 0, len(ts.tokens))
		for _, v := range ts.tokens {
			copyList = append(copyList, v)
		}
		ts.mu.RUnlock()

		f, err := os.Create(ts.tokensFile)
		if err != nil {
			return err
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(copyList)
	}

	// GenerateToken issues a new token. ttl is a time.Duration (0 = no expiry).
	func (ts *TokenStore) GenerateToken(username, name string, ttl time.Duration) (*Token, error) {
		tokStr, err := genToken(16)
		if err != nil {
			return nil, err
		}
		var exp *time.Time
		if ttl > 0 {
			t := time.Now().Add(ttl)
			exp = &t
		}
		tr := &Token{Token: tokStr, User: username, Name: name, ExpiresAt: exp, Revoked: false, CreatedAt: time.Now()}
		ts.mu.Lock()
		ts.tokens[tr.Token] = tr
		ts.mu.Unlock()
		log.Printf("TokenStore: generated token for user=%s name=%s token=%s", username, name, tr.Token)
		_ = ts.persistList()
		return tr, nil
	}

	// ListTokens returns all tokens.
	func (ts *TokenStore) ListTokens() []*Token {
		ts.mu.RLock()
		defer ts.mu.RUnlock()
		out := make([]*Token, 0, len(ts.tokens))
		for _, t := range ts.tokens {
			out = append(out, t)
		}
		return out
	}

	// Validate checks a token and returns it.
	func (ts *TokenStore) Validate(token string) (*Token, error) {
		ts.mu.RLock()
		defer ts.mu.RUnlock()
		t, ok := ts.tokens[token]
		if !ok || t.Revoked {
			return nil, errors.New("token not found")
		}
		if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
			return nil, errors.New("token expired")
		}
		log.Printf("TokenStore: validate: token ok: %s", token)
		return t, nil
	}

	// Rotate replaces an existing token with a new one.
	func (ts *TokenStore) Rotate(oldToken string, ttl time.Duration) (*Token, error) {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		tr, ok := ts.tokens[oldToken]
		if !ok {
			return nil, errors.New("token not found")
		}
		newTok, err := genToken(16)
		if err != nil {
			return nil, err
		}
		var exp *time.Time
		if ttl > 0 {
			t := time.Now().Add(ttl)
			exp = &t
		}
		newRec := &Token{Token: newTok, User: tr.User, Name: tr.Name, ExpiresAt: exp, Revoked: false, CreatedAt: time.Now()}
		ts.tokens[newTok] = newRec
		// revoke old token
		tr.Revoked = true
		delete(ts.tokens, oldToken)
		log.Printf("TokenStore: rotated token %s -> %s", oldToken, newTok)
		_ = ts.persistList()
		return newRec, nil
	}

	// Revoke marks a token revoked and persists.
	func (ts *TokenStore) Revoke(token string) error {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		if t, ok := ts.tokens[token]; ok {
			t.Revoked = true
			delete(ts.tokens, token)
			log.Printf("TokenStore: revoked token %s", token)
			return ts.persistList()
		}
		return nil
	}

	// StartWatcher starts a background goroutine to periodically reload tokens from disk.
	func (ts *TokenStore) StartWatcher() {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					_ = ts.load()
				case <-ts.stopChan:
					return
				}
			}
		}()
	}

	// StopWatcher signals the background watcher to stop.
	func (ts *TokenStore) StopWatcher() {
		select {
		case <-ts.stopChan:
			return
		default:
			close(ts.stopChan)
		}
	}

	// genToken produces a hex token string of n bytes
	func genToken(n int) (string, error) {
		if n <= 0 {
			n = 16
		}
		b := make([]byte, n)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return hex.EncodeToString(b), nil
	}
		return nil, err
	}
	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}
	tr := &Token{Token: tokStr, User: username, Name: name, ExpiresAt: exp, Revoked: false}
	ts.mu.Lock()
	ts.tokens[tr.Token] = tr
	ts.mu.Unlock()
	log.Printf("TokenStore: generated token for user=%s name=%s token=%s", username, name, tr.Token)
	_ = ts.persistList()
	return tr, nil
}

// ListTokens returns all tokens.
func (ts *TokenStore) ListTokens() []*Token {
	ts.mu.RLock()
	defer ts.mu.RUnlock()


		ts.mu.RLock()
		defer ts.mu.RUnlock()
		out := make([]*Token, 0, len(ts.tokens))
		for _, t := range ts.tokens {
			out = append(out, t)
		}
		return out
	}

	// Validate checks a token and returns it.
	func (ts *TokenStore) Validate(token string) (*Token, error) {
		ts.mu.RLock()
		defer ts.mu.RUnlock()
		t, ok := ts.tokens[token]
		if !ok || t.Revoked {
			return nil, errors.New("token not found")
		}
		if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
			return nil, errors.New("token expired")
		}
		log.Printf("TokenStore: validate: token ok: %s", token)
		return t, nil
	}

	// Rotate replaces an existing token with a new one.
	func (ts *TokenStore) Rotate(oldToken string, ttl time.Duration) (*Token, error) {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		tr, ok := ts.tokens[oldToken]
		if !ok {
			return nil, errors.New("token not found")
		}
		newTok, err := genToken(16)
		if err != nil {
			return nil, err
		}
		var exp *time.Time
		if ttl > 0 {
			t := time.Now().Add(ttl)
			exp = &t
		}
		newRec := &Token{Token: newTok, User: tr.User, Name: tr.Name, ExpiresAt: exp, Revoked: false}
		ts.tokens[newTok] = newRec
		// revoke old token
		tr.Revoked = true
		delete(ts.tokens, oldToken)
		log.Printf("TokenStore: rotated token %s -> %s", oldToken, newTok)
		_ = ts.persistList()
		return newRec, nil
	}

	// Revoke marks a token revoked and persists.
	func (ts *TokenStore) Revoke(token string) error {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		if t, ok := ts.tokens[token]; ok {
			t.Revoked = true
			delete(ts.tokens, token)
			log.Printf("TokenStore: revoked token %s", token)
			return ts.persistList()
		}
		return nil
	}

	// StartWatcher starts a background goroutine to periodically reload tokens from disk.
	func (ts *TokenStore) StartWatcher() {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					_ = ts.load()
				case <-ts.stopChan:
					return
				}
			}
		}()
	}

	// StopWatcher signals the background watcher to stop.
	func (ts *TokenStore) StopWatcher() {
		close(ts.stopChan)
	}

	}
	log.Printf("TokenStore: validate: token ok: %s", token)
	return tr, nil
}

func (ts *TokenStore) Rotate(oldToken string, ttl time.Duration) (*TokenRecord, error) {
	ts.mu.Lock()
	tr, ok := ts.tokens[oldToken]
	if !ok {
		ts.mu.Unlock()
		log.Printf("TokenStore: rotate: old token not found: %s", oldToken)
		return nil, errors.New("token not found")
	}
	// generate new
	newTok, err := genToken(16)
	if err != nil {
		return nil, err
	}
	var exp *time.Time
	if ttlSeconds > 0 {
		t := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
		exp = &t
	}
	newRec := &TokenRecord{Token: newTok, Name: tr.Name, Username: tr.Username, User: tr.User, ExpiresAt: exp}
	ts.tokens[newTok] = newRec
	delete(ts.tokens, oldToken)
	// make a snapshot and release lock before persisting
	list := make([]*TokenRecord, 0, len(ts.tokens))
	for _, v := range ts.tokens {
		list = append(list, v)
	}
	ts.mu.Unlock()
	log.Printf("TokenStore: rotated token %s -> %s", oldToken, newTok)
	if err := ts.persistList(list); err != nil {
		log.Printf("TokenStore: persist error after rotate: %v", err)
	}
	return newRec, nil
}

func (ts *TokenStore) Revoke(token string) error {
	ts.mu.Lock()
	if _, ok := ts.tokens[token]; !ok {
		ts.mu.Unlock()
		return nil
	}
	delete(ts.tokens, token)
	// snapshot and unlock before persisting
	list := make([]*TokenRecord, 0, len(ts.tokens))
	for _, v := range ts.tokens {
		list = append(list, v)
	}
	ts.mu.Unlock()
	log.Printf("TokenStore: revoked token %s", token)
	if err := ts.persistList(list); err != nil {
		log.Printf("TokenStore: persist error after revoke: %v", err)
		return err
	}
	return nil
}
