package bert

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBytesImageReader(t *testing.T) {
	const want = "image bytes"

	reader := newBytesImageReader([]byte(want))
	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if string(got) != want {
		t.Errorf("Read() = %q, want %q", got, want)
	}

	if err := reader.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestFileImageReaderRead(t *testing.T) {
	t.Parallel()

	const want = "image bytes"
	path := filepath.Join(t.TempDir(), "image.jpg")
	if err := os.WriteFile(path, []byte(want), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	reader := newFileImageReader(path)
	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if string(got) != want {
		t.Errorf("Read() = %q, want %q", got, want)
	}

	if err := reader.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestFileImageReaderReadMissingFile(t *testing.T) {
	reader := newFileImageReader(filepath.Join(t.TempDir(), "missing.jpg"))

	buf := make([]byte, 1)
	n, err := reader.Read(buf)
	if n != -1 {
		t.Errorf("Read() count = %d, want -1", n)
	}
	if err == nil {
		t.Fatal("Read() error = nil, want an error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Read() error = %v, want an os.ErrNotExist error", err)
	}
}

func TestFileImageReaderCloseBeforeRead(t *testing.T) {
	reader := newFileImageReader(filepath.Join(t.TempDir(), "image.jpg"))

	if err := reader.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}
