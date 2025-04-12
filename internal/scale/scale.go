package scale

import (
	"image"
	"log"
	"math"

	"golang.org/x/image/draw"
)

const (
	FocalTop    = "top"
	FocalCenter = "center"
	FocalBottom = "bottom"
)

// ScaleTo defines the target size and focal point for scaling an image.
// If the height is set to -1, the image will be scaled to the width only, (i.e., auto height).
type ScaleTo struct {
	Width  int
	Height int
	// Focal determines where the crop will be taken from.
	// Options are: scale.FocalTop, scale.FocalCenter, scale.FocalBottom
	// The default is scale.FocalTop.
	Focal string
}

func ScaleImage(src image.Image, scaleTo ScaleTo) image.Image {
	srcBounds := src.Bounds()
	scale := getScale(scaleTo, srcBounds.Dx(), srcBounds.Dy())
	scaledImage := image.NewRGBA(scale)
	draw.BiLinear.Scale(scaledImage, scaledImage.Rect, src, srcBounds, draw.Over, nil)

	if scaleTo.Height > 0 {
		crop := getCrop(scaleTo, scale.Dx(), scale.Dy())
		return scaledImage.SubImage(crop).(*image.RGBA)
	} else {
		return scaledImage
	}
}

func getScale(scaleTo ScaleTo, width int, height int) image.Rectangle {
	aspectRatio := float64(height) / float64(width)
	scaleToHeight := int(math.Round(float64(scaleTo.Width) * aspectRatio))
	return image.Rect(0, 0, scaleTo.Width, scaleToHeight)
}

func getCrop(scaleTo ScaleTo, width int, height int) image.Rectangle {
	cropHeight := min(scaleTo.Height, height)
	cropY := 0 // Default to top crop

	switch scaleTo.Focal {
	case FocalCenter:
		cropY = (height - cropHeight) / 2
	case FocalBottom:
		cropY = height - cropHeight
	}

	log.Printf("Crop Y: %d, Crop Height: %d", cropY, cropHeight)

	return image.Rect(0, cropY, width, cropY+cropHeight)
}
