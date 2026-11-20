package bert

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileOutputHandler(t *testing.T) {
	t.Run("directory", func(t *testing.T) {
		dir := t.TempDir()

		handler, ok := NewFileOutputHandler(dir).Value()
		if !ok {
			t.Fatal("NewFileOutputHandler() did not return a handler")
		}

		if handler == nil {
			t.Fatal("NewFileOutputHandler() returned a nil handler")
		}
	})

	t.Run("file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "output.jpg")
		if err := os.WriteFile(path, []byte("output"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		_, ok := NewFileOutputHandler(path).Value()
		if ok {
			t.Fatal("NewFileOutputHandler() returned a handler for a file")
		}
	})

	t.Run("missing path", func(t *testing.T) {
		_, ok := NewFileOutputHandler(filepath.Join(t.TempDir(), "missing")).Value()
		if ok {
			t.Fatal("NewFileOutputHandler() returned a handler for a missing path")
		}
	})
}

func TestFileOutputHandlerHandlePNG(t *testing.T) {
	dir := t.TempDir()
	handler, ok := NewFileOutputHandler(dir).Value()
	if !ok {
		t.Fatal("NewFileOutputHandler() did not return a handler")
	}

	const name = "output.png"
	img := testOutputImage()
	result := handler.Handle(processedImage{
		image: img,
		name:  name,
		transform: Transform{
			OutputType: PngOutputType,
		},
	})
	if _, ok := result.Value(); !ok {
		t.Fatal("Handle() did not return a value")
	}

	file, err := os.Open(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer file.Close()

	got, err := png.Decode(file)
	if err != nil {
		t.Fatalf("png.Decode() error = %v", err)
	}
	if got.Bounds() != img.Bounds() {
		t.Errorf("decoded bounds = %v, want %v", got.Bounds(), img.Bounds())
	}
}

func TestFileOutputHandlerHandleJPEG(t *testing.T) {
	dir := t.TempDir()
	handler, ok := NewFileOutputHandler(dir).Value()
	if !ok {
		t.Fatal("NewFileOutputHandler() did not return a handler")
	}

	const name = "output.jpg"
	result := handler.Handle(processedImage{
		image: testOutputImage(),
		name:  name,
		transform: Transform{
			OutputType: JpgOutputType,
			Quality:    FullQuality,
		},
	})
	if _, ok := result.Value(); !ok {
		t.Fatal("Handle() did not return a value")
	}

	file, err := os.Open(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer file.Close()

	got, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("jpeg.Decode() error = %v", err)
	}
	if got.Bounds() != (image.Rect(0, 0, 2, 2)) {
		t.Errorf("decoded bounds = %v, want 2x2", got.Bounds())
	}
}

func TestFileOutputHandlerHandleCreateError(t *testing.T) {
	dir := t.TempDir()
	handler, ok := NewFileOutputHandler(dir).Value()
	if !ok {
		t.Fatal("NewFileOutputHandler() did not return a handler")
	}

	_, ok = handler.Handle(processedImage{
		image: testOutputImage(),
		name:  filepath.Join("missing", "output.jpg"),
		transform: Transform{
			OutputType: JpgOutputType,
		},
	}).Value()
	if ok {
		t.Fatal("Handle() returned a value when output creation failed")
	}
}

func TestTransformGetJpegQuality(t *testing.T) {
	tests := []struct {
		name    string
		quality string
		want    int
	}{
		{name: "full", quality: FullQuality, want: 100},
		{name: "high", quality: HighQuality, want: 90},
		{name: "medium", quality: MediumQuality, want: 75},
		{name: "low", quality: LowQuality, want: 60},
		{name: "default", quality: "unknown", want: 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Transform{Quality: tt.quality}).GetJpegQuality(); got != tt.want {
				t.Errorf("GetJpegQuality() = %d, want %d", got, tt.want)
			}
		})
	}
}

func testOutputImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	return img
}
