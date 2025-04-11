package scale

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

type scale struct {
	width  int
	height int
}

func ScaleImage(src image.Image, width int) *image.RGBA {
	inputBounds := src.Bounds()
	scale := getScale(inputBounds.Dx(), inputBounds.Dy(), width)
	scaledBounds := getScaleBounds(scale)
	scaledImage := image.NewRGBA(scaledBounds)
	draw.BiLinear.Scale(scaledImage, scaledImage.Rect, src, inputBounds, draw.Over, nil)
	return scaledImage
}

func getScaleBounds(scale *scale) image.Rectangle {
	return image.Rect(0, 0, scale.width, scale.height)
}

func getScale(width int, height int, scaleWidth int) *scale {
	aspectRatio := float64(height) / float64(width)
	scaleHeight := int(math.Round(float64(scaleWidth) * aspectRatio))
	return &scale{
		width:  scaleWidth,
		height: scaleHeight}
}
