package bert

import (
	"image"
	img "image"
	"io"
	"path/filepath"
	"strings"

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

type CropWindow struct {
	Left, Top, Right, Bottom int
	Dimensions
}

type UnprocessedImage struct {
	Source io.Reader
	Name   string
	// OutputType string // jpg, png
}

type ImageTransform struct {
	Image     img.Image
	Name      string
	Transform Transform
}

type ProcessedImage struct {
	Image     img.Image
	Name      string
	Transform Transform
}

type ImageSource interface {
	Next() option.Option[UnprocessedImage]
}

type OutputHandler interface {
	Handle(ProcessedImage) option.Option[error]
}

type ImageTransformer interface {
	Transform(ImageTransform) result.Result[ProcessedImage, error]
}

// Check if the file is an image by checking the file extension
func IsImage(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case JpgExtension, JpegExtension, PngExtension, WebpExtension:
		return true
	default:
		return false
	}
}

func GetDimensions(img image.Image) Dimensions {
	bnds := img.Bounds()
	return Dimensions{Width: bnds.Dx(), Height: bnds.Dy()}
}
