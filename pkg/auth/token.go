package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/yourusername/ops-tool/pkg/audit"
)

// Token represents an API token
type Token struct {
	Token     string     `json:"token"`
	User      string     `json:"user"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Revoked   bool       `json:"revoked"`
}

// TokenStore manages tokens persisted to disk
type TokenStore struct {
	mu     sync.RWMutex
	tokens map[string]*Token
	file   string
}

// NewTokenStore creates a token store (default file: tokens.json)
func NewTokenStore(file string) *TokenStore {
	if file == "" {
		file = "tokens.json"
	}
	ts := &TokenStore{
		tokens: make(map[string]*Token),
		file:   file,
	}
	ts.load()
	return ts
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
	_ = audit.Record("token.create", user, tok.Token, map[string]any{"name": name})
	return tok, nil
}

// Validate checks a token and returns the associated token object
func (s *TokenStore) Validate(t string) (*Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tok, ok := s.tokens[t]
	if !ok {
		return nil, errors.New("token not found")
	}
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
		_ = audit.Record("token.revoke", tok.User, tok.Token, nil)
	}
	return err
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
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(s.tokens)
}

// load loads tokens from disk
func (s *TokenStore) load() {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.Open(s.file)
	if err != nil {
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	_ = dec.Decode(&s.tokens)
}
