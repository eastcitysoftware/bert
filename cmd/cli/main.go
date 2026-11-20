package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eastcitysoftware/bert/internal/bert"
)

const (
	version = "0.0.1"
)

func main() {
	printBanner(os.Stdout)
	if result, err := run(); err != nil {
		printError(os.Stderr, err)
		os.Exit(1)
	} else {
		printResult(os.Stdout, result)
	}
}

func run() (*bert.TransformResult, error) {
	inputPath := flag.String("input", "", "input path, can be a file or directory")
	outputPath := flag.String("output", "", "the output directory path")
	mode := flag.String("mode", "scale", "resize mode. Options: scale, crop")
	outputType := flag.String("type", "jpg", "image output type. Options: jpg, png")
	quality := flag.String("quality", "high", "quality of the output image (jpg only). Options: full, high, medium and low")
	width := flag.Int("width", 640, "desired output width")
	height := flag.Int("height", -1, "desired output height. If > -1, image is cropped to this height")
	cropx := flag.String("crop-x", "center", "crop position. Options: left, right, center")
	cropy := flag.String("crop-y", "center", "crop position. Options: top, bottom, center")

	flag.Parse()

	// was the program called with no arguments?
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	// Use only the non-flag arguments as positional input and output paths.
	applyPositionalPaths(inputPath, outputPath, flag.Args())

	// verify input
	if *inputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("input path does not exist: %s", *inputPath)
	}

	if *outputPath == "" {
		return nil, fmt.Errorf("output path is required")
	}

	if err := ensureOutputPath(*outputPath); err != nil {
		return nil, err
	}

	fileImageSource := bert.NewFileImageSource(*inputPath, *outputType)
	fileImageErr, isErr := fileImageSource.ErrorValue()

	if isErr {
		return nil, fileImageErr
	}

	imageSource, isOk := fileImageSource.Value()

	if !isOk {
		return nil, fmt.Errorf("An unknown error occurred when accessing the image source")
	}

	outputHandlerResult := bert.NewFileOutputHandler(*outputPath)
	outputHandlerErr, isErr := outputHandlerResult.ErrorValue()

	if isErr {
		return nil, outputHandlerErr
	}

	outputHandler, isOk := outputHandlerResult.Value()

	if !isOk {
		return nil, fmt.Errorf("An unknown error occurred when accessing the output handler")
	}

	result := bert.RunTransform(bert.TransformOptions{
		ImageSource:   imageSource,
		OutputHandler: outputHandler,
		Transform: bert.Transform{
			Mode:       *mode,
			OutputType: *outputType,
			Quality:    *quality,
			Dimensions: bert.Dimensions{
				Width:  *width,
				Height: *height},
			Crop: bert.Crop{
				CropX: *cropx,
				CropY: *cropy,
			},
		},
	})

	if e, isErr := result.ErrorValue(); isErr {
		return nil, e
	}

	if v, isOk := result.Value(); isOk {
		return &v, nil
	}

	return nil, fmt.Errorf("An unknown error occurred when accessing the transform result")
}

func applyPositionalPaths(inputPath, outputPath *string, positionalArgs []string) {
	if len(positionalArgs) >= 1 && *inputPath == "" {
		*inputPath = positionalArgs[0]
	}

	if len(positionalArgs) >= 2 && *outputPath == "" {
		*outputPath = positionalArgs[1]
	}
}

func ensureOutputPath(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("unable to create output directory %q: %w", path, err)
	}

	outputStat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("unable to access output directory %q: %w", path, err)
	}
	if !outputStat.IsDir() {
		return fmt.Errorf("output path is not a directory: %s", path)
	}

	return nil
}

func printResult(w io.Writer, r *bert.TransformResult) {
	trns := fmt.Sprintf("mode: %s / output: %s (%s) / dimensions: %v",
		r.Transform.Mode,
		r.Transform.OutputType,
		r.Transform.Quality,
		r.Transform.Dimensions)

	var files []string

	for i, result := range r.Results {
		files = append(files, fmt.Sprintf("#%3d, %s", i+1, result.Name))
	}

	fmt.Fprintf(w, `
* Processing took: %s

* Transform: %s

* File(s):
  %s

`, r.Duration, trns, strings.Join(files, "\n  "))
}

// source: https://patorjk.com/software/taag/#p=display&f=Modular&t=BERT&x=none&v=4&h=4&w=75&we=false
func printBanner(w io.Writer) {
	fmt.Fprintf(w, `
 _______  _______  ______    _______
|  _    ||       ||    _ |  |       |
| |_|   ||    ___||   | ||  |_     _|
|       ||   |___ |   |_||_   |   |
|  _   | |    ___||    __  |  |   |
| |_|   ||   |___ |   |  | |  |   |
|_______||_______||___|  |_|  |___|   %s
  `, version)
}

func printError(w io.Writer, err error) {
	fmt.Fprintf(w, `
! ERROR
  %v

  `, err)
}
