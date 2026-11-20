package bert

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"testing"
)

func TestTransformValidator(t *testing.T) {
	tests := map[string]struct {
		transform Transform
		isOk      bool
	}{
		"invalid": {Transform{}, false},
		"valid":   {Transform{TransformModeScale, PngOutputType, HighQuality, Dimensions{200, 200}, Crop{CropXCenter, CropYCenter}}, true},
	}

	validator := newTransformValidator()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if validator.Validate(test.transform).IsOk() != test.isOk {
				t.Errorf("DimensionsValidator should be Ok, got Error")
			}
		})
	}
}

func TestDimensionsValidator(t *testing.T) {
	tests := map[string]struct {
		crop Dimensions
		isOk bool
	}{
		"invalid":   {Dimensions{}, false},
		"invalid h": {Dimensions{0, 0}, false},
		"invalid w": {Dimensions{-1, 0}, false},
		"valid w":   {Dimensions{200, 0}, true},
		"valid h":   {Dimensions{0, 200}, true},
		"valid":     {Dimensions{200, 200}, true},
	}

	validator := newDimensionsValidator()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if validator.Validate(test.crop).IsOk() != test.isOk {
				t.Errorf("DimensionsValidator should be Ok, got Error")
			}
		})
	}
}
func TestCropValidator(t *testing.T) {
	tests := map[string]struct {
		crop Crop
		isOk bool
	}{
		"invalid":       {Crop{}, false},
		"center/center": {Crop{CropXCenter, CropYCenter}, true},
		"center/top":    {Crop{CropXCenter, CropYTop}, true},
		"center/bottom": {Crop{CropXCenter, CropYBottom}, true},
		"left/center":   {Crop{CropXLeft, CropYCenter}, true},
		"left/top":      {Crop{CropXLeft, CropYTop}, true},
		"left/bottom":   {Crop{CropXLeft, CropYBottom}, true},
		"right/center":  {Crop{CropXRight, CropYCenter}, true},
		"right/top":     {Crop{CropXRight, CropYTop}, true},
		"right/bottom":  {Crop{CropXRight, CropYBottom}, true},
	}

	validator := newCropValidator()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if validator.Validate(test.crop).IsOk() != test.isOk {
				t.Errorf("CropValidator should be Ok, got Error")
			}
		})
	}
}

func TestUnprocessedImageDecode(t *testing.T) {
	pngData := encodePNG(t)
	jpegData := encodeJPEG(t)

	tests := []struct {
		name       string
		source     []byte
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "test.png",
			source:     pngData,
			wantWidth:  2,
			wantHeight: 1,
		},
		{
			name:       "test.jpeg",
			source:     jpegData,
			wantWidth:  2,
			wantHeight: 1,
		},
		{
			name:       "test.jpg",
			source:     jpegData,
			wantWidth:  2,
			wantHeight: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := (unprocessedImage{
				source: func() ImageReader { return newBytesImageReader(tt.source) },
				name:   tt.name,
			}).decode().Value()
			if !ok {
				t.Fatal("Decode() returned an error")
			}

			bounds := got.Bounds()
			if (tt.wantWidth > 0 && bounds.Dx() != tt.wantWidth) ||
				(tt.wantHeight > 0 && bounds.Dy() != tt.wantHeight) ||
				bounds.Dx() <= 0 || bounds.Dy() <= 0 {
				t.Errorf("Decode() bounds = %v, want %dx%d", bounds, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}

func TestUnprocessedImageDecodeErrors(t *testing.T) {
	tests := []struct {
		name    string
		source  []byte
		wantErr string
	}{
		{
			name:    "image.gif",
			wantErr: "unsupported image format: .gif",
		},
		{
			name:    "malformed.png",
			source:  []byte("not an image"),
			wantErr: "png: invalid format: not a PNG file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, ok := (unprocessedImage{
				source: func() ImageReader { return newBytesImageReader(tt.source) },
				name:   tt.name,
			}).decode().ErrorValue()
			if !ok {
				t.Fatal("Decode() did not return an error")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Decode() error = %q, want %q", err, tt.wantErr)
			}
		})
	}
}

func decodeBase64(t *testing.T, data string) []byte {
	t.Helper()

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		t.Fatalf("decode base64 fixture: %v", err)
	}
	return decoded
}

func encodePNG(t *testing.T) []byte {
	t.Helper()

	return encodeImage(t, png.Encode)
}

func encodeJPEG(t *testing.T) []byte {
	t.Helper()

	return encodeImage(t, func(w io.Writer, img image.Image) error {
		return jpeg.Encode(w, img, nil)
	})
}

func encodeImage(t *testing.T, encode func(io.Writer, image.Image) error) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{B: 255, A: 255})

	var data bytes.Buffer
	if err := encode(&data, img); err != nil {
		t.Fatalf("encode image: %v", err)
	}
	return data.Bytes()
}
