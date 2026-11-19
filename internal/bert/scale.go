package bert

import (
	"image"
	"math"

	"github.com/eastcitysoftware/danom/result"
	"golang.org/x/image/draw"
)

type ScaleToWidth struct {
	Dimensions
	OutputType string // jpg, png, webp
	Quality    string // full, high, medium, low
	CropY      string // top, bottom, center
}

type ScaleToHeight struct {
	Dimensions
	OutputType string // jpg, png, webp
	Quality    string // full, high, medium, low
	CropX      string // left, right, center
}

func TransformToWidth(transform Transform) ScaleToWidth {
	return ScaleToWidth{
		Dimensions: transform.Dimensions,
		OutputType: transform.OutputType,
		Quality:    transform.Quality,
		CropY:      transform.CropY,
	}
}

func TransformToHeight(transform Transform) ScaleToHeight {
	return ScaleToHeight{
		Dimensions: transform.Dimensions,
		OutputType: transform.OutputType,
		Quality:    transform.Quality,
		CropX:      transform.CropX,
	}
}

func (scale ScaleToWidth) Transform(src ImageTransform) result.Result[ProcessedImage, error] {
	scaleDim := getScale(scale.Dimensions, GetDimensions(src.Image))
	return src.scaleTo(scaleDim)
}

func (scale ScaleToHeight) Transform(src ImageTransform) result.Result[ProcessedImage, error] {
	scaleDim := getScale(scale.Dimensions, GetDimensions(src.Image))
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

func (src ImageTransform) scaleTo(toDim Dimensions) result.Result[ProcessedImage, error] {
	scale := image.Rect(0, 0, toDim.Width, toDim.Height)
	scaledImage := image.NewRGBA(scale)

	draw.ApproxBiLinear.Scale(
		scaledImage,
		scaledImage.Rect,
		src.Image,
		src.Image.Bounds(),
		draw.Over,
		nil)

	return result.Ok[error](ProcessedImage{
		Image:     scaledImage,
		Name:      src.Name,
		Transform: src.Transform})
}
