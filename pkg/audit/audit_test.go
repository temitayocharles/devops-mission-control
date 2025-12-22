package audit

import (
	"os"
	"testing"
)

func TestRecordAndRead(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// record
	if err := Record("", "test.action", "bob", "target", map[string]any{"k": "v"}); err != nil {
		t.Fatal(err)
	}
	entries, err := ReadEntries("")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != "test.action" {
		t.Fatalf("unexpected action: %s", entries[0].Action)
	}
	// cleanup file
	_ = os.Remove(DefaultFile)
}
