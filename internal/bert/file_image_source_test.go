package bert

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestIsImage(t *testing.T) {
	if isImage("some-file.txt") {
		t.Errorf("IsImage() should be false for .txt")
	}

	tests := map[string]struct {
		Extension string
	}{
		JpgExtension:  {Extension: JpgExtension},
		JpegExtension: {Extension: JpegExtension},
		PngExtension:  {Extension: PngExtension},
		WebpExtension: {Extension: WebpExtension},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			if !isImage("some-file" + test.Extension) {
				t.Errorf("IsImage() should be true for %s", test.Extension)
			}
		})
	}
}

func TestGetImages(t *testing.T) {
	t.Run("directory filters image files", func(t *testing.T) {
		dir := t.TempDir()
		writeSourceTestFile(t, filepath.Join(dir, "first.jpg"), "first")
		writeSourceTestFile(t, filepath.Join(dir, "ignored.txt"), "ignored")
		writeSourceTestFile(t, filepath.Join(dir, "second.PNG"), "second")

		got, err := getImages(dir).Value()
		if !err {
			t.Fatal("getImages() did not return a value")
		}

		want := []string{
			filepath.Join(dir, "first.jpg"),
			filepath.Join(dir, "second.PNG"),
		}
		if !slices.Equal(got, want) {
			t.Errorf("getImages() = %v, want %v", got, want)
		}
	})

	t.Run("single image file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "image.webp")
		writeSourceTestFile(t, path, "image")

		got, ok := getImages(path).Value()
		if !ok {
			t.Fatal("getImages() did not return a value")
		}

		want := []string{path}
		if !slices.Equal(got, want) {
			t.Errorf("getImages() = %v, want %v", got, want)
		}
	})

	t.Run("non-image file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "notes.txt")
		writeSourceTestFile(t, path, "notes")

		got, ok := getImages(path).Value()
		if !ok {
			t.Fatal("getImages() did not return a value")
		}

		if len(got) != 0 {
			t.Errorf("getImages() = %v, want no files", got)
		}
	})

	t.Run("missing path", func(t *testing.T) {
		_, ok := getImages(filepath.Join(t.TempDir(), "missing.jpg")).Value()
		if ok {
			t.Fatal("getImages() returned a value, want an error")
		}

		err, ok := getImages(filepath.Join(t.TempDir(), "missing.jpg")).ErrorValue()
		if !ok {
			t.Fatal("getImages() did not return an error")
		}
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("getImages() error = %v, want an os.ErrNotExist error", err)
		}
	})
}

func TestNewFileImageSource(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "image.jpg")
	writeSourceTestFile(t, path, "image")

	for _, outputType := range []string{JpgOutputType, PngOutputType} {
		t.Run(outputType, func(t *testing.T) {
			source, ok := NewFileImageSource(path, outputType).Value()
			if !ok {
				t.Fatal("NewFileImageSource() did not return a source")
			}

			next := source.Next()
			if next.IsNone() {
				t.Fatal("Next() returned none for an existing image")
			}

			image, ok := next.Value()
			if !ok {
				t.Fatal("Next() did not return an image")
			}
			if image.name != filepath.Base(path) {
				t.Errorf("Next() name = %q, want %q", image.name, filepath.Base(path))
			}

			if source.Next().IsSome() {
				t.Error("Next() returned an image after the source was exhausted")
			}
		})
	}

	t.Run("invalid output type", func(t *testing.T) {
		_, ok := NewFileImageSource(path, "gif").Value()
		if ok {
			t.Fatal("NewFileImageSource() returned a source for an invalid output type")
		}
	})
}

func TestFileImageSourceNextSkipsMissingFile(t *testing.T) {
	source := &fileImageSource{
		filePaths: []string{filepath.Join(t.TempDir(), "missing.jpg")},
	}

	if source.Next().IsSome() {
		t.Fatal("Next() returned an image for a missing file")
	}
	if source.Next().IsSome() {
		t.Fatal("Next() returned an image after the missing file was processed")
	}
}

func TestGetUnprocessedImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "image.jpeg")
	const want = "image data"
	writeSourceTestFile(t, path, want)

	image, ok := getUnprocessedImage(path, PngOutputType).Value()
	if !ok {
		t.Fatal("getUnprocessedImage() did not return an image")
	}
	if image.name != filepath.Base(path) {
		t.Errorf("Name = %q, want %q", image.name, filepath.Base(path))
	}

	reader := image.source()
	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Source().Read() error = %v", err)
	}
	if string(got) != want {
		t.Errorf("Source().Read() = %q, want %q", got, want)
	}
	if err := reader.Close(); err != nil {
		t.Errorf("Source().Close() error = %v", err)
	}

	if getUnprocessedImage(filepath.Join(t.TempDir(), "missing.jpg"), JpgOutputType).IsSome() {
		t.Error("getUnprocessedImage() returned an image for a missing file")
	}
}

func writeSourceTestFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
