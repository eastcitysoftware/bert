package bert

import (
	"image"

	"github.com/eastcitysoftware/danom/result"
)

type CropToDimensions struct {
	Dimensions
	Crop
	OutputType string // jpg, png, webp
	Quality    string // full, high, medium, low
}

func TransformToDimensions(t Transform) CropToDimensions {
	return CropToDimensions{
		Dimensions: t.Dimensions,
		OutputType: t.OutputType,
		Crop:       t.Crop,
	}
}

func (crop CropToDimensions) Transform(src ImageTransform) result.Result[ProcessedImage, error] {
	cropWin := getCropWindow(crop.Crop, src.Transform.Dimensions, GetDimensions(src.Image))
	return src.cropTo(cropWin)
}

func getCropWindow(crop Crop, dim, srcDim Dimensions) CropWindow {
	x := 0
	switch crop.CropX {
	case CropXCenter:
		x = (srcDim.Width - dim.Width) / 2
	case CropXRight:
		x = srcDim.Width - dim.Width
	}

	y := 0
	switch crop.CropY {
	case CropYCenter:
		y = (srcDim.Height - dim.Height) / 2
	case CropYBottom:
		y = srcDim.Height - dim.Height
	}
	return CropWindow{
		Left:       x,
		Top:        y,
		Right:      dim.Width + x,
		Bottom:     dim.Height + y,
		Dimensions: dim,
	}
}

func (src ImageTransform) cropTo(window CropWindow) result.Result[ProcessedImage, error] {
	croppedImage := src.Image.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(window.Left, window.Top, window.Right, window.Bottom))

	return result.Ok[error](ProcessedImage{
		Image:     croppedImage,
		Name:      src.Name,
		Transform: src.Transform})
}
