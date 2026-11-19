package bert

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/eastcitysoftware/danom/option"
	"github.com/eastcitysoftware/danom/result"
)

type FileOutputHandler struct {
	OutputPath string
}

func NewFileOutputHandler(outputPath string) result.Result[OutputHandler, error] {
	return result.Bind(result.Of(os.Stat(outputPath)),
		func(f os.FileInfo) result.Result[OutputHandler, error] {
			if !f.IsDir() {
				return result.Error[OutputHandler](fmt.Errorf("Output path must be a directory"))
			}
			return result.Ok[error, OutputHandler](&FileOutputHandler{outputPath})
		})
}

func (out *FileOutputHandler) Handle(img ProcessedImage) option.Option[error] {
	outputPath := filepath.Join(out.OutputPath, img.Name)

	transformed, err := os.Create(outputPath)
	if err != nil {
		return option.Some(fmt.Errorf("failed to create resize file: %v", err))
	}
	defer transformed.Close()

	switch img.Transform.OutputType {
	case PngOutputType:
		png.Encode(transformed, img.Image)
	default:
		jpeg.Encode(transformed, img.Image, &jpeg.Options{Quality: img.Transform.GetJpegQuality()})
	}
	return option.None[error]()
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
