package auth

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "errors"
    "os"
    "sync"
    "time"
)

// User represents a local user for RBAC
type User struct {
    Username  string    `json:"username"`
    Role      string    `json:"role"`
    Token     string    `json:"token"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
}

// GetUser returns the User with the given username or an error.
func (u *UserStore) GetUser(username string) (*User, error) {
    users, err := LoadUsers()
    if err != nil {
        return nil, err
    }
    _, user := FindUser(users, username)
    if user == nil {
        return nil, os.ErrNotExist
    }
    return user, nil
}

// ListUsers returns all users.
func (u *UserStore) ListUsers() []User {
    users, _ := LoadUsers()
    return users
}

// DeleteUser removes a user.
func (u *UserStore) DeleteUser(username string) error {
    return RemoveUser(username)
}

// SetUserRole sets the role for a user.
func (u *UserStore) SetUserRole(username string, role Role) error {
    return SetRole(username, string(role))
}

// Authenticate returns the user if credentials are valid.
// Note: password storage is not implemented in this simple store; we treat
// existence as successful authentication for now.
func (u *UserStore) Authenticate(username, _password string) (*User, error) {
    users, err := LoadUsers()
    if err != nil {
        return nil, err
    }
    _, user := FindUser(users, username)
    if user == nil {
        return nil, os.ErrNotExist
    }
    return user, nil
}

var (
    usersFile = "users.json"
    mu        sync.RWMutex
)

func genToken(n int) (string, error) {
    if n <= 0 {
        n = 16
    }
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

// LoadUsers reads users from users.json
func LoadUsers() ([]User, error) {
    mu.RLock()
    defer mu.RUnlock()
    f, err := os.Open(usersFile)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return []User{}, nil
        }
        return nil, err
    }
    defer f.Close()
    var u []User
    dec := json.NewDecoder(f)
    if err := dec.Decode(&u); err != nil {
        return nil, err
    }
    return u, nil
}

// SaveUsers writes users to users.json
func SaveUsers(users []User) error {
    mu.Lock()
    defer mu.Unlock()
    f, err := os.Create(usersFile)
    if err != nil {
        return err
    }
    defer f.Close()
    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")
    return enc.Encode(users)
}

// FindUser finds a user by username
func FindUser(users []User, username string) (int, *User) {
    for i := range users {
        if users[i].Username == username {
            return i, &users[i]
        }
    }
    return -1, nil
}

// AddUser creates a user and returns a token
func AddUser(username, role string) (string, error) {
    users, err := LoadUsers()
    if err != nil {
        return "", err
    }
    if _, u := FindUser(users, username); u != nil {
        return "", nil
    }
    tok, err := genToken(16)
    if err != nil {
        return "", err
    }
    u := User{Username: username, Role: role, Token: tok, CreatedAt: time.Now()}
    users = append(users, u)
    if err := SaveUsers(users); err != nil {
        return "", err
    }
    return tok, nil
}

// RemoveUser deletes a user
func RemoveUser(username string) error {
    users, err := LoadUsers()
    if err != nil {
        return err
    }
    idx, _ := FindUser(users, username)
    if idx == -1 {
        return nil
    }
    users = append(users[:idx], users[idx+1:]...)
    return SaveUsers(users)
}

// SetRole updates role for a user
func SetRole(username, role string) error {
    users, err := LoadUsers()
    if err != nil {
        return err
    }
    idx, _ := FindUser(users, username)
    if idx == -1 {
        return nil
    }
    users[idx].Role = role
    return SaveUsers(users)
}

// AuthenticateToken checks username+token
func AuthenticateToken(username, token string) bool {
    users, err := LoadUsers()
    if err != nil {
        return false
    }
    _, u := FindUser(users, username)
    if u == nil {
        return false
    }
    return u.Token == token
}
