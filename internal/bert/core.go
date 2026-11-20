package bert

import (
	"fmt"
	"image"
	"io"

	"github.com/eastcitysoftware/danom/option"
	"github.com/eastcitysoftware/danom/result"
)

const (
	TransformModeCrop  = "crop"
	TransformModeScale = "scale"

	JpgOutputType = "jpg"
	PngOutputType = "png"

	JpegExtension = ".jpeg"
	JpgExtension  = ".jpg"
	PngExtension  = ".png"
	WebpExtension = ".webp"

	FullQuality   = "full"
	HighQuality   = "high"
	MediumQuality = "medium"
	LowQuality    = "low"

	CropYTop    = "top"
	CropYCenter = "center"
	CropYBottom = "bottom"

	CropXLeft   = "left"
	CropXRight  = "right"
	CropXCenter = "center"
)

type Transform struct {
	Mode       string // scale, crop
	OutputType string // jpg, png
	Quality    string // full, high, medium, low
	Dimensions
	Crop
}

type Dimensions struct {
	Width  int
	Height int
}

type Crop struct {
	CropX string // left, right, center
	CropY string // top, bottom, center
}

type unprocessedImage struct {
	source func() ImageReader
	name   string
}

type imageTransform struct {
	image     image.Image
	name      string
	transform Transform
}

type processedImage struct {
	image      image.Image
	name       string
	dimensions Dimensions
	transform  Transform
}

type unit struct{}

type ImageReader interface {
	io.Reader
	io.Closer
}

type ImageSource interface {
	Next() option.Option[unprocessedImage]
}

type OutputHandler interface {
	Handle(processedImage) result.Result[unit, error]
}

type imageTransformer interface {
	transform(imageTransform) result.Result[processedImage, error]
}

func getDimensions(img image.Image) Dimensions {
	bnds := img.Bounds()
	return Dimensions{Width: bnds.Dx(), Height: bnds.Dy()}
}

func (d Dimensions) String() string {
	switch {
	case d.Height > 0 && d.Width > 0:
		return fmt.Sprintf("%dW x %dH", d.Width, d.Height)
	case d.Height > 0:
		return fmt.Sprintf("Auto x %dH", d.Height)
	default:
		return fmt.Sprintf("%dW x Auto", d.Width)
	}
}
