package bert

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/eastcitysoftware/danom/option"
	"github.com/eastcitysoftware/danom/result"
)

type fileImageSource struct {
	outputType     string
	filePaths      []string
	processedPaths []string
}

func NewFileImageSource(inputPath string, outputType string) result.Result[ImageSource, error] {
	switch outputType {
	case JpgOutputType, PngOutputType:
		return result.Map(getImages(inputPath), func(files []string) ImageSource {
			return &fileImageSource{filePaths: files}
		})
	default:
		return result.Error[ImageSource](fmt.Errorf("Invalid output type, must be one of: %s, %s", JpgOutputType, PngOutputType))
	}
}

func (src *fileImageSource) Next() option.Option[unprocessedImage] {
	if len(src.filePaths) == len(src.processedPaths) {
		return option.None[unprocessedImage]()
	}

	for _, f := range src.filePaths {
		if slices.Contains(src.processedPaths, f) {
			continue
		}

		src.processedPaths = append(src.processedPaths, f)
		return getUnprocessedImage(f, src.outputType)
	}
	return option.None[unprocessedImage]()
}

func getUnprocessedImage(inputPath, outputType string) option.Option[unprocessedImage] {
	_, err := os.Stat(inputPath)

	if err != nil {
		return option.None[unprocessedImage]()
	}

	return option.Some(unprocessedImage{
		source: func() ImageReader { return newFileImageReader(inputPath) },
		name:   filepath.Base(inputPath)})
}

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
			if isImage(inputPath) {
				files = append(files, inputPath)
			}
		}
	} else {
		if isImage(inputPath) {
			files = []string{inputPath}
		}
	}

	return result.Ok[error](files)
}

// Check if the file is an image by checking the file extension
func isImage(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case JpgExtension, JpegExtension, PngExtension, WebpExtension:
		return true
	default:
		return false
	}
}
