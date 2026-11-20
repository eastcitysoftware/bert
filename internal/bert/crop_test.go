package bert

import (
	"image"
	"image/color"
	"testing"
)

func TestGetCropWindow(t *testing.T) {
	tests := map[string]struct {
		crop     Crop
		dim      Dimensions
		srcDim   Dimensions
		expected cropWindow
	}{
		"left/top": {
			crop:   Crop{CropX: CropXLeft, CropY: CropYTop},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 0, top: 0, right: 200, bottom: 100,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
		"center/top": {
			crop:   Crop{CropX: CropXCenter, CropY: CropYTop},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 100, top: 0, right: 300, bottom: 100,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
		"right/top": {
			crop:   Crop{CropX: CropXRight, CropY: CropYTop},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 200, top: 0, right: 400, bottom: 100,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
		"left/center": {
			crop:   Crop{CropX: CropXLeft, CropY: CropYCenter},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 0, top: 100, right: 200, bottom: 200,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
		"center/center": {
			crop:   Crop{CropX: CropXCenter, CropY: CropYCenter},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 100, top: 100, right: 300, bottom: 200,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
		"right/center": {
			crop:   Crop{CropX: CropXRight, CropY: CropYCenter},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 200, top: 100, right: 400, bottom: 200,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
		"left/bottom": {
			crop:   Crop{CropX: CropXLeft, CropY: CropYBottom},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 0, top: 200, right: 200, bottom: 300,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
		"center/bottom": {
			crop:   Crop{CropX: CropXCenter, CropY: CropYBottom},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 100, top: 200, right: 300, bottom: 300,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
		"right/bottom": {
			crop:   Crop{CropX: CropXRight, CropY: CropYBottom},
			dim:    Dimensions{Width: 200, Height: 100},
			srcDim: Dimensions{Width: 400, Height: 300},
			expected: cropWindow{
				left: 200, top: 200, right: 400, bottom: 300,
				Dimensions: Dimensions{Width: 200, Height: 100},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := getCropWindow(test.crop, test.dim, test.srcDim)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestImageTransformCropTo(t *testing.T) {
	srcImage := image.NewRGBA(image.Rect(0, 0, 4, 4))
	colors := map[image.Point]color.RGBA{
		image.Pt(1, 1): {R: 255, A: 255},
		image.Pt(2, 1): {G: 255, A: 255},
		image.Pt(1, 2): {B: 255, A: 255},
		image.Pt(2, 2): {R: 255, G: 255, A: 255},
	}
	for point, c := range colors {
		srcImage.SetRGBA(point.X, point.Y, c)
	}

	got, ok := (imageTransform{image: srcImage, name: "test.png"}).cropTo(cropWindow{
		left: 1, top: 1, right: 3, bottom: 3,
		Dimensions: Dimensions{Width: 2, Height: 2},
	}).Value()
	if !ok {
		t.Fatal("cropTo() returned an error")
	}

	wantBounds := image.Rect(1, 1, 3, 3)
	if got.image.Bounds() != wantBounds {
		t.Fatalf("cropTo() bounds = %v, want %v", got.image.Bounds(), wantBounds)
	}

	tests := map[string]struct {
		point image.Point
		want  color.RGBA
	}{
		"top-left":     {image.Pt(1, 1), colors[image.Pt(1, 1)]},
		"top-right":    {image.Pt(2, 1), colors[image.Pt(2, 1)]},
		"bottom-left":  {image.Pt(1, 2), colors[image.Pt(1, 2)]},
		"bottom-right": {image.Pt(2, 2), colors[image.Pt(2, 2)]},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			gotR, gotG, gotB, gotA := got.image.At(test.point.X, test.point.Y).RGBA()
			wantR, wantG, wantB, wantA := test.want.RGBA()
			if gotR != wantR || gotG != wantG || gotB != wantB || gotA != wantA {
				t.Errorf("pixel at %v = (%d, %d, %d, %d), want (%d, %d, %d, %d)",
					test.point, gotR, gotG, gotB, gotA,
					wantR, wantG, wantB, wantA)
			}
		})
	}
}

func TestCropToDimensionsTransform(t *testing.T) {
	srcImage := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			srcImage.SetRGBA(x, y, color.RGBA{
				R: uint8(10 + x*40),
				G: uint8(20 + y*40),
				B: 255,
				A: 255,
			})
		}
	}

	const cropWidth = 2
	const cropHeight = 2
	tests := map[string]struct {
		crop     Crop
		wantLeft int
		wantTop  int
	}{
		"left/top":      {Crop{CropX: CropXLeft, CropY: CropYTop}, 0, 0},
		"center/top":    {Crop{CropX: CropXCenter, CropY: CropYTop}, 1, 0},
		"right/top":     {Crop{CropX: CropXRight, CropY: CropYTop}, 2, 0},
		"left/center":   {Crop{CropX: CropXLeft, CropY: CropYCenter}, 0, 1},
		"center/center": {Crop{CropX: CropXCenter, CropY: CropYCenter}, 1, 1},
		"right/center":  {Crop{CropX: CropXRight, CropY: CropYCenter}, 2, 1},
		"left/bottom":   {Crop{CropX: CropXLeft, CropY: CropYBottom}, 0, 2},
		"center/bottom": {Crop{CropX: CropXCenter, CropY: CropYBottom}, 1, 2},
		"right/bottom":  {Crop{CropX: CropXRight, CropY: CropYBottom}, 2, 2},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			src := imageTransform{
				image: srcImage,
				name:  "test.png",
				transform: Transform{
					Dimensions: Dimensions{Width: cropWidth, Height: cropHeight},
				},
			}
			crop := cropToDimensions{
				Dimensions: Dimensions{Width: cropWidth, Height: cropHeight},
				Crop:       test.crop,
			}

			got, ok := crop.transform(src).Value()
			if !ok {
				t.Fatal("Transform() returned an error")
			}

			wantBounds := image.Rect(
				test.wantLeft,
				test.wantTop,
				test.wantLeft+cropWidth,
				test.wantTop+cropHeight,
			)
			if got.image.Bounds() != wantBounds {
				t.Fatalf("Transform() bounds = %v, want %v", got.image.Bounds(), wantBounds)
			}

			for y := wantBounds.Min.Y; y < wantBounds.Max.Y; y++ {
				for x := wantBounds.Min.X; x < wantBounds.Max.X; x++ {
					if got.image.At(x, y) != srcImage.At(x, y) {
						t.Errorf("pixel at (%d, %d) = %v, want %v",
							x, y, got.image.At(x, y), srcImage.At(x, y))
					}
				}
			}
		})
	}
}
