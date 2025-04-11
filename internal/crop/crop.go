package crop

import "image"

type crop struct {
	width  int
	height int
	y      int
}

func CropImage(scaledImage *image.RGBA, width int, height int, cropHeight int, focalPoint string) image.Image {
	crop := getCrop(width, height, cropHeight, focalPoint)
	cropBounds := getCropBounds(crop)
	return scaledImage.SubImage(cropBounds)
}

func getCropBounds(crop *crop) image.Rectangle {
	return image.Rect(0, crop.y, crop.width, crop.y+crop.height)
}

func getCrop(width int, height int, cropHeight int, focalPoint string) *crop {
	cropY := 0 // Default to top crop
	switch focalPoint {
	case "center":
		cropY = (height - cropHeight) / 2
	case "bottom":
		cropY = height - cropHeight
	}

	return &crop{
		width:  width,
		height: cropHeight,
		y:      cropY}
}
