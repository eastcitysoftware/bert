package bert

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/eastcitysoftware/danom/result"
)

type fileOutputHandler struct {
	outputPath string
}

func NewFileOutputHandler(outputPath string) result.Result[OutputHandler, error] {
	return result.Bind(result.Of(os.Stat(outputPath)),
		func(f os.FileInfo) result.Result[OutputHandler, error] {
			if !f.IsDir() {
				return result.Error[OutputHandler](fmt.Errorf("Output path must be a directory"))
			}
			return result.Ok[error, OutputHandler](&fileOutputHandler{outputPath})
		})
}

func (out *fileOutputHandler) Handle(img processedImage) result.Result[unit, error] {
	outputPath := filepath.Join(out.outputPath, img.name)

	transformed, err := os.Create(outputPath)
	if err != nil {
		return result.Error[unit](fmt.Errorf("failed to create resize file: %v", err))
	}
	defer transformed.Close()

	switch img.transform.OutputType {
	case PngOutputType:
		png.Encode(transformed, img.image)
	default:
		jpeg.Encode(transformed, img.image, &jpeg.Options{Quality: img.transform.GetJpegQuality()})
	}
	return result.Ok[error](unit{})
}

func (t Transform) GetJpegQuality() int {
	switch t.Quality {
	case FullQuality:
		return 100
	case HighQuality:
		return 90
	case MediumQuality:
		return 75
	default:
		return 60
	}
}
