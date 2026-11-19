package bert

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/eastcitysoftware/danom/option"
	"github.com/eastcitysoftware/danom/result"
)

type FileImageSource struct {
	OutputType     string
	FilePaths      []string
	ProcessedPaths []string
}

func NewFileImageSource(inputPath string, outputType string) result.Result[ImageSource, error] {
	switch outputType {
	case JpgOutputType, PngOutputType:
		return result.Map(getImages(inputPath), func(files []string) ImageSource {
			return &FileImageSource{FilePaths: files}
		})
	default:
		return result.Error[ImageSource](fmt.Errorf("Invalid output type, must be one of: %s, %s", JpgOutputType, PngOutputType))
	}
}

func (src *FileImageSource) Next() option.Option[UnprocessedImage] {
	if len(src.FilePaths) == len(src.ProcessedPaths) {
		return option.None[UnprocessedImage]()
	}

	for _, f := range src.FilePaths {
		if slices.Contains(src.ProcessedPaths, f) {
			continue
		}

		src.ProcessedPaths = append(src.ProcessedPaths, f)
		return getUnprocessedImage(f, src.OutputType)
	}
	return option.None[UnprocessedImage]()
}

func getUnprocessedImage(inputPath, outputType string) option.Option[UnprocessedImage] {
	r, err := os.Open(inputPath)
	if err != nil {
		return option.None[UnprocessedImage]()
	}
	name := filepath.Base(inputPath)
	return option.Some(UnprocessedImage{
		Source: r,
		Name:   name}) //,
	// OutputType: outputType})
}

// func getOutputType(inputPath string) string {
// 	switch filepath.Ext(inputPath) {
// 	case JpegExtension, JpgExtension:
// 		return JpgOutputType
// 	case PngExtension
// 	}
// }

func getImages(inputPath string) result.Result[[]string, error] {
	inputStat, err := os.Stat(inputPath)
	if err != nil {
		return result.Error[[]string](err)
	}

	var files []string
	if inputStat.IsDir() {
		dirFiles, err := os.ReadDir(inputPath)
		if err != nil {
			return result.Error[[]string](err)
		}
		for _, file := range dirFiles {
			inputPath := filepath.Join(inputPath, file.Name())
			if IsImage(inputPath) {
				files = append(files, inputPath)
			}
		}
	} else {
		if IsImage(inputPath) {
			files = []string{inputPath}
		}
	}

	return result.Ok[error](files)
}
