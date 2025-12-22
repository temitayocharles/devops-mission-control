package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents an audit log entry
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Action    string         `json:"action"`
	Actor     string         `json:"actor,omitempty"`
	Target    string         `json:"target,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
}

var (
	DefaultFile = "audit.log.jsonl"
	mu          sync.Mutex
)

// Record appends an audit entry to the audit file (JSON Lines)
func Record(action, actor, target string, details map[string]any) error {
	mu.Lock()
	defer mu.Unlock()
	e := Entry{
		Timestamp: time.Now(),
		Action:    action,
		Actor:     actor,
		Target:    target,
		Details:   details,
	}
	f, err := os.OpenFile(DefaultFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "audit file close failed: %v\n", cerr)
		}
	}()
	enc := json.NewEncoder(f)
	return enc.Encode(e)
}

// ReadEntries reads audit entries from the given file (JSON Lines). If file=="" uses DefaultFile.
func ReadEntries(file string) ([]Entry, error) {
	if file == "" {
		file = DefaultFile
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "audit file close failed: %v\n", cerr)
		}
	}()
	dec := json.NewDecoder(f)
	var res []Entry
	for {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}
