package crop

import "testing"

func TestGetCropBounds(t *testing.T) {
	crop := getCrop(1280, 960, 480, "top")
	bounds := getCropBounds(crop)
	if bounds.Min.X != 0 || bounds.Min.Y != 0 {
		t.Errorf("Expected bounds Min (0, 0), got (%d, %d)", bounds.Min.X, bounds.Min.Y)
	}
	if bounds.Max.X != 1280 || bounds.Max.Y != 480 {
		t.Errorf("Expected bounds Max (1280, 480), got (%d, %d)", bounds.Max.X, bounds.Max.Y)
	}
}

func TestGetCrop(t *testing.T) {
	crop := getCrop(1280, 960, 480, "top")
	if crop.width != 1280 {
		t.Errorf("Expected width 1280, got %d", crop.width)
	}
	if crop.height != 480 {
		t.Errorf("Expected height 480, got %d", crop.height)
	}
	if crop.y != 0 {
		t.Errorf("Expected y 0, got %d", crop.y)
	}

	crop = getCrop(1280, 960, 480, "center")
	if crop.y != 240 {
		t.Errorf("Expected y 240, got %d", crop.y)
	}

	crop = getCrop(1280, 960, 480, "bottom")
	if crop.y != 480 {
		t.Errorf("Expected y 480, got %d", crop.y)
	}
}
