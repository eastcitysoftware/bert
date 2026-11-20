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

	got := transform.newImageTransform("test.png", img)

	if got.image != img {
		t.Fatal("NewImageTransform() did not preserve the image")
	}
	if got.transform != transform {
		t.Errorf("NewImageTransform() transform = %#v, want %#v", got.transform, transform)
	}
}
