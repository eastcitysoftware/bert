package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/eastcitysoftware/bert/internal/bert"
)

const (
	version = "0.0.1"
)

func main() {
	printBanner(os.Stdout)
	if err := run(); err != nil {
		printError(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
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

	// was the program called with unnamed arguments?
	if len(os.Args) >= 2 {
		if *inputPath == "" {
			inputPath = &os.Args[1]
		}
	}

	if len(os.Args) >= 3 {
		if *outputPath == "" {
			outputPath = &os.Args[2]
		}
	}

	// verify input
	if *inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input path does not exist: %s", *inputPath)
	}

	if *outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	if err := ensureOutputPath(*outputPath); err != nil {
		return err
	}

	fileImageSource := bert.NewFileImageSource(*inputPath, *outputType)
	fileImageErr, isErr := fileImageSource.ErrorValue()

	if isErr {
		return fileImageErr
	}

	imageSource, isOk := fileImageSource.Value()

	if !isOk {
		return fmt.Errorf("An unknown error occurred when accessing the image source")
	}

	outputHandlerResult := bert.NewFileOutputHandler(*outputPath)
	outputHandlerErr, isErr := outputHandlerResult.ErrorValue()

	if isErr {
		return outputHandlerErr
	}

	outputHandler, isOk := outputHandlerResult.Value()

	if !isOk {
		return fmt.Errorf("An unknown error occurred when accessing the output handler")
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
		return e
	}

	return nil
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

// func processImage(src imageToProcess, resize input) error {
// 	original, err := os.Open(src.inputPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open src file: %v", err)
// 	}
// 	defer original.Close()

// 	resized, err := os.Create(src.outputPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create resize file: %v", err)
// 	}
// 	defer resized.Close()

// 	decodedImage, err := decodeImage(original, getMimetype(src.inputPath))
// 	if err != nil {
// 		return fmt.Errorf("failed to decode image: %v", err)
// 	}

// 	// Scale the image to the desired width while maintaining the aspect ratio
// 	scaledImage := scale.ScaleImage(decodedImage, scale.ScaleTo{
// 		Width:  resize.width,
// 		Height: resize.height,
// 		Focal:  resize.focalPoint})

// 	switch {
// 	case resize.outputMimetype == pngMimeType:
// 		return png.Encode(resized, scaledImage)
// 	default:
// 		return jpeg.Encode(resized, scaledImage, &jpeg.Options{Quality: resize.quality})
// 	}
// }

// func decodeImage(original io.Reader, mimeType string) (image.Image, error) {
// 	var src image.Image
// 	var err error

// 	switch mimeType {
// 	case jpegMimeType:
// 		src, err = jpeg.Decode(original)
// 	case pngMimeType:
// 		src, err = png.Decode(original)
// 	case webpMimeType:
// 		src, err = webp.Decode(original)
// 	default:
// 		src, err = nil, fmt.Errorf("unsupported image format: %s", mimeType)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return src, nil
// }

// func getImages(inputPath string, outputPath string, outputExtension string) ([]imageToProcess, error) {
// 	inputStat, err := os.Stat(inputPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var files []string
// 	if inputStat.IsDir() {
// 		dirFiles, err := os.ReadDir(inputPath)
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, file := range dirFiles {
// 			files = append(files, filepath.Join(inputPath, file.Name()))
// 		}
// 	} else {
// 		files = []string{inputPath}
// 	}

// 	var images []imageToProcess

// 	for _, file := range files {
// 		if bert.IsImage(file) {
// 			outputName := strings.Replace(filepath.Base(file), filepath.Ext(file), "."+outputExtension, 1)
// 			outputPath := filepath.Join(outputPath, outputName)

// 			images = append(images, imageToProcess{
// 				inputPath:  file,
// 				outputPath: outputPath})
// 		}
// 	}

// 	return images, nil
// }
