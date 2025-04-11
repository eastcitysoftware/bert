package scale

import "testing"

func TestGetScaleBounds(t *testing.T) {
	scale := getScale(1280, 960, 640)
	bounds := getScaleBounds(scale)
	if bounds.Min.X != 0 || bounds.Min.Y != 0 {
		t.Errorf("Expected bounds Min (0, 0), got (%d, %d)", bounds.Min.X, bounds.Min.Y)
	}
	if bounds.Max.X != 640 || bounds.Max.Y != 480 {
		t.Errorf("Expected bounds Max (640, 480), got (%d, %d)", bounds.Max.X, bounds.Max.Y)
	}
}

func TestGetScale(t *testing.T) {
	scale := getScale(1280, 960, 640)
	if scale.width != 640 {
		t.Errorf("Expected width 640, got %d", scale.width)
	}
	if scale.height != 480 {
		t.Errorf("Expected height 480, got %d", scale.height)
	}

	scale = getScale(1280, 960, 1280)
	if scale.width != 1280 {
		t.Errorf("Expected width 1280, got %d", scale.width)
	}
	if scale.height != 960 {
		t.Errorf("Expected height 960, got %d", scale.height)
	}
}
