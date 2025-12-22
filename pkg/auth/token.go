package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yourusername/ops-tool/pkg/audit"
)

// Token represents an API token
type Token struct {
	Token     string     `json:"token"`
	TenantID  string     `json:"tenant_id,omitempty"`
	User      string     `json:"user"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Revoked   bool       `json:"revoked"`
}

// TokenStore manages tokens persisted to disk
type TokenStore struct {
	mu      sync.RWMutex
	tokens  map[string]*Token
	file    string
	lastMod time.Time
	modMu   sync.Mutex
	stopCh  chan struct{}
	tenant  string
}

// NewTokenStore creates a token store (default file: tokens.json)
func NewTokenStore(tenant, file string) *TokenStore {
	if file == "" {
		if tenant == "" {
			file = "tokens.json"
		} else {
			file = fmt.Sprintf("tokens.%s.json", tenant)
		}
	}
	ts := &TokenStore{
		tokens: make(map[string]*Token),
		file:   file,
		tenant: tenant,
	}
	ts.load()
	return ts
}

// StartWatcher starts the background file-watcher that reloads the tokens file on changes.
func (s *TokenStore) StartWatcher() {
	if s.stopCh != nil {
		return
	}
	s.stopCh = make(chan struct{})
	// start file watcher (fsnotify preferred)
	go s.watchFile()
	// start background sweeper to revoke expired tokens
	go s.sweeper()
}

// StopWatcher stops the background file-watcher.
func (s *TokenStore) StopWatcher() {
	if s.stopCh != nil {
		close(s.stopCh)
		s.stopCh = nil
	}
}

func genToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateToken creates and stores a new token for a user
func (s *TokenStore) GenerateToken(user, name string, ttl time.Duration) (*Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tstr, err := genToken(32)
	if err != nil {
		return nil, err
	}
	var exp *time.Time
	if ttl > 0 {
		tm := time.Now().Add(ttl)
		exp = &tm
	}
	tok := &Token{
		Token:     tstr,
		TenantID:  s.tenant,
		User:      user,
		Name:      name,
		CreatedAt: time.Now(),
		ExpiresAt: exp,
		Revoked:   false,
	}
	s.tokens[tstr] = tok
	if err := s.save(); err != nil {
		return nil, err
	}
	if err := audit.Record(s.tenant, "token.create", user, tok.Token, map[string]any{"name": name}); err != nil {
		fmt.Fprintf(os.Stderr, "audit record failed: %v\n", err)
	}
	return tok, nil
}

// Validate checks a token and returns the associated token object
func (s *TokenStore) Validate(t string) (*Token, error) {
	// Check in-memory first.
	s.mu.RLock()
	tok, ok := s.tokens[t]
	s.mu.RUnlock()

	if !ok {
		// not present in memory — attempt to reload from disk and re-check
		s.load()
		s.mu.RLock()
		tok, ok = s.tokens[t]
		s.mu.RUnlock()
		if !ok {
			return nil, errors.New("token not found")
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if tok.Revoked {
		return nil, errors.New("token revoked")
	}
	if tok.ExpiresAt != nil && time.Now().After(*tok.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return tok, nil
}

// Revoke marks a token revoked
func (s *TokenStore) Revoke(t string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tok, ok := s.tokens[t]
	if !ok {
		return errors.New("token not found")
	}
	tok.Revoked = true
	err := s.save()
	if err == nil {
		if rerr := audit.Record(s.tenant, "token.revoke", tok.User, tok.Token, nil); rerr != nil {
			fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
		}
	}
	return err
}

// Rotate replaces an existing token with a new one for the same user/name.
// The old token is revoked. Returns the new token.
func (s *TokenStore) Rotate(old string, ttl time.Duration) (*Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tok, ok := s.tokens[old]
	if !ok {
		return nil, errors.New("token not found")
	}
	if tok.Revoked {
		return nil, errors.New("token revoked")
	}
	// generate new token
	tstr, err := genToken(32)
	if err != nil {
		return nil, err
	}
	var exp *time.Time
	if ttl > 0 {
		tm := time.Now().Add(ttl)
		exp = &tm
	}
	newTok := &Token{
		Token:     tstr,
		User:      tok.User,
		Name:      tok.Name,
		CreatedAt: time.Now(),
		ExpiresAt: exp,
		Revoked:   false,
	}
	// revoke old and store new
	tok.Revoked = true
	s.tokens[newTok.Token] = newTok
	if err := s.save(); err != nil {
		return nil, err
	}
	if rerr := audit.Record(s.tenant, "token.rotate", newTok.User, newTok.Token, map[string]any{"replaced": old}); rerr != nil {
		fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
	}
	return newTok, nil
}

// ListTokens returns all tokens
func (s *TokenStore) ListTokens() []*Token {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*Token, 0, len(s.tokens))
	for _, t := range s.tokens {
		res = append(res, t)
	}
	return res
}

// save persists tokens to disk
func (s *TokenStore) save() error {
	f, err := os.Create(s.file)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "tokens file close failed: %v\n", cerr)
		}
	}()
	enc := json.NewEncoder(f)
	if err := enc.Encode(s.tokens); err != nil {
		return err
	}
	// update lastMod
	if fi, err := os.Stat(s.file); err == nil {
		s.modMu.Lock()
		s.lastMod = fi.ModTime()
		s.modMu.Unlock()
	}
	return nil
}

// load loads tokens from disk
func (s *TokenStore) load() {
	s.mu.Lock()
	defer s.mu.Unlock()
	fi, err := os.Stat(s.file)
	if err != nil {
		return
	}
	f, err := os.Open(s.file)
	if err != nil {
		return
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "tokens file close failed: %v\n", cerr)
		}
	}()
	dec := json.NewDecoder(f)
	_ = dec.Decode(&s.tokens)
	s.modMu.Lock()
	s.lastMod = fi.ModTime()
	s.modMu.Unlock()
}

// watchFile uses fsnotify when available; falls back to polling if watcher fails.
func (s *TokenStore) watchFile() {
	// attempt fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		defer func() {
			if cerr := watcher.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "watcher close failed: %v\n", cerr)
			}
		}()
		dir := filepath.Dir(s.file)
		// watch the directory so we see creates/renames
		if err := watcher.Add(dir); err == nil {
			for {
				select {
				case ev, ok := <-watcher.Events:
					if !ok {
						return
					}
					// if the event is for our file and it's a write/create/rename, reload
					if filepath.Clean(ev.Name) == filepath.Clean(s.file) {
						if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
							s.load()
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					if rerr := audit.Record(s.tenant, "tokenwatch.error", "", s.file, map[string]any{"error": err.Error()}); rerr != nil {
						fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
					}
				case <-s.stopCh:
					return
				}
			}
		}
	}
	// fallback: polling
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fi, err := os.Stat(s.file)
			if err != nil {
				continue
			}
			s.modMu.Lock()
			lm := s.lastMod
			s.modMu.Unlock()
			if fi.ModTime().After(lm) {
				s.load()
			}
		case <-s.stopCh:
			return
		}
	}
}

// sweeper periodically revokes expired tokens and persists changes.
func (s *TokenStore) sweeper() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			changed := false
			s.mu.Lock()
			for k, t := range s.tokens {
				if t.Revoked {
					continue
				}
				if t.ExpiresAt != nil && now.After(*t.ExpiresAt) {
					t.Revoked = true
					if rerr := audit.Record(s.tenant, "token.expire.revoke", t.User, t.Token, nil); rerr != nil {
						fmt.Fprintf(os.Stderr, "audit record failed: %v\n", rerr)
					}
					changed = true
				}
				// safety: remove entries that are revoked and older than 30 days to keep file small
				if t.Revoked && t.CreatedAt.Add(30*24*time.Hour).Before(now) {
					delete(s.tokens, k)
					changed = true
				}
			}
			if changed {
				_ = s.save()
			}
			s.mu.Unlock()
		case <-s.stopCh:
			return
		}
	}
}
