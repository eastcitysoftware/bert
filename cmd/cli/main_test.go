package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureOutputPath(t *testing.T) {
	t.Run("creates missing directory and parents", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "parent", "output")

		if err := ensureOutputPath(path); err != nil {
			t.Fatalf("ensureOutputPath() error = %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat output path: %v", err)
		}
		if !info.IsDir() {
			t.Fatal("output path is not a directory")
		}
	})

	t.Run("rejects file path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "output")
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}

		if err := ensureOutputPath(path); err == nil {
			t.Fatal("ensureOutputPath() expected an error")
		}
	})
}
