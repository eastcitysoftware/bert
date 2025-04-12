package scale

import (
	"image"
	"testing"
)

func TestScaleImage(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1280, 960))
	scaleTo := ScaleTo{Width: 640, Height: 200, Focal: FocalCenter}
	scaledImage := ScaleImage(src, scaleTo)

	if scaledImage.Bounds().Dx() != 640 || scaledImage.Bounds().Dy() != 200 {
		t.Errorf("Expected scaled image size (640, 200), got (%d, %d)", scaledImage.Bounds().Dx(), scaledImage.Bounds().Dy())
	}
}

func TestGetScale(t *testing.T) {
	scale := getScale(ScaleTo{Width: 640, Focal: FocalCenter}, 1280, 960)
	if scale.Dx() != 640 || scale.Dy() != 480 {
		t.Errorf("Expected scale (640, 480), got (%d, %d)", scale.Dx(), scale.Dy())
	}
	if scale.Min.X != 0 || scale.Min.Y != 0 {
		t.Errorf("Expected scale Min (0, 0), got (%d, %d)", scale.Min.X, scale.Min.Y)
	}
	if scale.Max.X != 640 || scale.Max.Y != 480 {
		t.Errorf("Expected scale Max (640, 480), got (%d, %d)", scale.Max.X, scale.Max.Y)
	}
}

func TestGetCrop(t *testing.T) {
	crop := getCrop(ScaleTo{Width: 640, Height: 200, Focal: FocalCenter}, 640, 480)
	if crop.Min.X != 0 || crop.Min.Y != 140 {
		t.Errorf("Expected crop Min (0, 140), got (%d, %d)", crop.Min.X, crop.Min.Y)
	}
	if crop.Max.X != 640 || crop.Max.Y != 340 {
		t.Errorf("Expected crop Max (640, 340), got (%d, %d)", crop.Max.X, crop.Max.Y)
	}
	if crop.Dx() != 640 || crop.Dy() != 200 {
		t.Errorf("Expected crop (640, 200), got (%d, %d)", crop.Dx(), crop.Dy())
	}

	crop = getCrop(ScaleTo{Width: 640, Height: 200, Focal: FocalTop}, 640, 480)
	if crop.Min.X != 0 || crop.Min.Y != 0 {
		t.Errorf("Expected crop Min (0, 0), got (%d, %d)", crop.Min.X, crop.Min.Y)
	}
	if crop.Max.X != 640 || crop.Max.Y != 200 {
		t.Errorf("Expected crop Max (640, 200), got (%d, %d)", crop.Max.X, crop.Max.Y)
	}
	if crop.Dx() != 640 || crop.Dy() != 200 {
		t.Errorf("Expected crop (640, 200), got (%d, %d)", crop.Dx(), crop.Dy())
	}

	crop = getCrop(ScaleTo{Width: 640, Height: 200, Focal: FocalBottom}, 640, 480)
	if crop.Min.X != 0 || crop.Min.Y != 280 {
		t.Errorf("Expected crop Min (0, 280), got (%d, %d)", crop.Min.X, crop.Min.Y)
	}
	if crop.Max.X != 640 || crop.Max.Y != 480 {
		t.Errorf("Expected crop Max (640, 480), got (%d, %d)", crop.Max.X, crop.Max.Y)
	}
	if crop.Dx() != 640 || crop.Dy() != 200 {
		t.Errorf("Expected crop (640, 200), got (%d, %d)", crop.Dx(), crop.Dy())
	}
}
