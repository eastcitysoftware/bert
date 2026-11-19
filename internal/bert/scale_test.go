package bert

import (
	"image"
	"image/color"
	"testing"
)

func TestGetScale(t *testing.T) {
	tests := map[string]struct {
		dim      Dimensions
		srcDim   Dimensions
		expected Dimensions
	}{
		"scale by width, square": {
			dim:      Dimensions{Width: 200, Height: 0},
			srcDim:   Dimensions{Width: 400, Height: 400},
			expected: Dimensions{Width: 200, Height: 200},
		},
		"scale by height, square": {
			dim:      Dimensions{Width: 0, Height: 200},
			srcDim:   Dimensions{Width: 400, Height: 400},
			expected: Dimensions{Width: 200, Height: 200},
		},
		"scale by width, irregular": {
			dim:      Dimensions{Width: 200, Height: 0},
			srcDim:   Dimensions{Width: 400, Height: 300},
			expected: Dimensions{Width: 200, Height: 150},
		},
		"scale by height, irregular": {
			dim:      Dimensions{Width: 0, Height: 200},
			srcDim:   Dimensions{Width: 400, Height: 300},
			expected: Dimensions{Width: 267, Height: 200},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := getScale(test.dim, test.srcDim)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestImageTransformScaleTo(t *testing.T) {
	srcImage := image.NewRGBA(image.Rect(0, 0, 2, 2))
	srcColor := color.RGBA{R: 255, G: 64, B: 128, A: 255}
	for y := range 2 {
		for x := range 2 {
			srcImage.SetRGBA(x, y, srcColor)
		}
	}

	got, ok := (ImageTransform{
		Image: srcImage,
		Name:  "test.png",
	}).scaleTo(Dimensions{Width: 4, Height: 3}).Value()
	if !ok {
		t.Fatal("scaleTo() returned an error")
	}

	wantBounds := image.Rect(0, 0, 4, 3)
	if got.Image.Bounds() != wantBounds {
		t.Fatalf("scaleTo() bounds = %v, want %v", got.Image.Bounds(), wantBounds)
	}

	wantR, wantG, wantB, wantA := srcColor.RGBA()
	for y := wantBounds.Min.Y; y < wantBounds.Max.Y; y++ {
		for x := wantBounds.Min.X; x < wantBounds.Max.X; x++ {
			gotR, gotG, gotB, gotA := got.Image.At(x, y).RGBA()
			if gotR != wantR || gotG != wantG || gotB != wantB || gotA != wantA {
				t.Errorf("pixel at (%d, %d) = (%d, %d, %d, %d), want (%d, %d, %d, %d)",
					x, y, gotR, gotG, gotB, gotA,
					wantR, wantG, wantB, wantA)
			}
		}
	}
}

func TestScaleToWidthTransform(t *testing.T) {
	testScaleTransform(t,
		ScaleToWidth{Dimensions: Dimensions{Width: 8}},
		Dimensions{Width: 8, Height: 4},
	)
}

func TestScaleToHeightTransform(t *testing.T) {
	testScaleTransform(t,
		ScaleToHeight{Dimensions: Dimensions{Height: 4}},
		Dimensions{Width: 8, Height: 4},
	)
}

func testScaleTransform(t *testing.T, transformer ImageTransformer, want Dimensions) {
	t.Helper()

	srcImage := image.NewRGBA(image.Rect(0, 0, 4, 2))
	srcColor := color.RGBA{R: 255, G: 64, B: 128, A: 255}
	for y := range 2 {
		for x := range 4 {
			srcImage.SetRGBA(x, y, srcColor)
		}
	}

	got, ok := transformer.Transform(ImageTransform{
		Image: srcImage,
		Name:  "test.png",
	}).Value()
	if !ok {
		t.Fatal("Transform() returned an error")
	}

	wantBounds := image.Rect(0, 0, want.Width, want.Height)
	if got.Image.Bounds() != wantBounds {
		t.Fatalf("Transform() bounds = %v, want %v", got.Image.Bounds(), wantBounds)
	}

	wantR, wantG, wantB, wantA := srcColor.RGBA()
	for y := wantBounds.Min.Y; y < wantBounds.Max.Y; y++ {
		for x := wantBounds.Min.X; x < wantBounds.Max.X; x++ {
			gotR, gotG, gotB, gotA := got.Image.At(x, y).RGBA()
			if gotR != wantR || gotG != wantG || gotB != wantB || gotA != wantA {
				t.Errorf("pixel at (%d, %d) = (%d, %d, %d, %d), want (%d, %d, %d, %d)",
					x, y, gotR, gotG, gotB, gotA,
					wantR, wantG, wantB, wantA)
			}
		}
	}
}
