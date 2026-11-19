package bert

import (
	"image"
	"testing"
)

func TestTransformNewImageTransform(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	transform := Transform{
		Mode:       TransformModeCrop,
		OutputType: PngOutputType,
		Quality:    HighQuality,
		Dimensions: Dimensions{Width: 100, Height: 50},
		Crop:       Crop{CropX: CropXCenter, CropY: CropYBottom},
	}

	got := transform.NewImageTransform("test.png", img)

	if got.Image != img {
		t.Fatal("NewImageTransform() did not preserve the image")
	}
	if got.Transform != transform {
		t.Errorf("NewImageTransform() transform = %#v, want %#v", got.Transform, transform)
	}
}

func TestIsImage(t *testing.T) {
	if IsImage("some-file.txt") {
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
			if !IsImage("some-file" + test.Extension) {
				t.Errorf("IsImage() should be true for %s", test.Extension)
			}
		})
	}
}
