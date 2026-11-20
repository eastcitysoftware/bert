package bert

import (
	"image"

	"github.com/eastcitysoftware/danom/result"
)

type cropWindow struct {
	left, top, right, bottom int
	Dimensions
}

type cropToDimensions struct {
	Dimensions
	Crop
	outputType string // jpg, png, webp
	quality    string // full, high, medium, low
}

func transformToDimensions(t Transform) cropToDimensions {
	return cropToDimensions{
		Dimensions: t.Dimensions,
		outputType: t.OutputType,
		Crop:       t.Crop,
	}
}

func (crop cropToDimensions) transform(src imageTransform) result.Result[processedImage, error] {
	cropWin := getCropWindow(crop.Crop, src.transform.Dimensions, getDimensions(src.image))
	return src.cropTo(cropWin)
}

func getCropWindow(crop Crop, dim, srcDim Dimensions) cropWindow {
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
	return cropWindow{
		left:       x,
		top:        y,
		right:      dim.Width + x,
		bottom:     dim.Height + y,
		Dimensions: dim,
	}
}

func (src imageTransform) cropTo(window cropWindow) result.Result[processedImage, error] {
	croppedImage := src.image.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(window.left, window.top, window.right, window.bottom))

	return result.Ok[error](processedImage{
		image:     croppedImage,
		name:      src.name,
		transform: src.transform})
}
