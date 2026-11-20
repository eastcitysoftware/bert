package bert

import (
	"image"
	"math"

	"github.com/eastcitysoftware/danom/result"
	"golang.org/x/image/draw"
)

type scaleToWidth struct {
	Dimensions
	outputType string // jpg, png, webp
	quality    string // full, high, medium, low
	cropY      string // top, bottom, center
}

type scaleToHeight struct {
	Dimensions
	outputType string // jpg, png, webp
	quality    string // full, high, medium, low
	cropX      string // left, right, center
}

func transformToWidth(transform Transform) scaleToWidth {
	return scaleToWidth{
		Dimensions: transform.Dimensions,
		outputType: transform.OutputType,
		quality:    transform.Quality,
		cropY:      transform.CropY,
	}
}

func transformToHeight(transform Transform) scaleToHeight {
	return scaleToHeight{
		Dimensions: transform.Dimensions,
		outputType: transform.OutputType,
		quality:    transform.Quality,
		cropX:      transform.CropX,
	}
}

func (scale scaleToWidth) transform(src imageTransform) result.Result[processedImage, error] {
	scaleDim := getScale(scale.Dimensions, getDimensions(src.image))
	return src.scaleTo(scaleDim)
}

func (scale scaleToHeight) transform(src imageTransform) result.Result[processedImage, error] {
	scaleDim := getScale(scale.Dimensions, getDimensions(src.image))
	return src.scaleTo(scaleDim)
}

func getScale(dim Dimensions, srcDim Dimensions) Dimensions {
	var scaleToWidth, scaleToHeight int
	if dim.Height > 0 {
		scaleToHeight = dim.Height
		scaleToWidth = int(math.Round(float64(dim.Height) * (float64(srcDim.Width) / float64(srcDim.Height))))
	} else {
		scaleToWidth = dim.Width
		scaleToHeight = int(math.Round(float64(dim.Width) * (float64(srcDim.Height) / float64(srcDim.Width))))
	}

	return Dimensions{
		Width:  scaleToWidth,
		Height: scaleToHeight}
}

func (src imageTransform) scaleTo(toDim Dimensions) result.Result[processedImage, error] {
	scale := image.Rect(0, 0, toDim.Width, toDim.Height)
	scaledImage := image.NewRGBA(scale)

	draw.ApproxBiLinear.Scale(
		scaledImage,
		scaledImage.Rect,
		src.image,
		src.image.Bounds(),
		draw.Over,
		nil)



	return result.Ok[error](processedImage{
		image:     scaledImage,
		name:      src.name,
		transform: src.transform})
}
