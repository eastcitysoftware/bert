package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

type inputFlags struct {
	inputPath  *string // absolute path to the input file or directory
	outputPath *string // absolute path to the output directory
	width      *int    // required, desired output width
	height     *int    // if > -1, image is cropped to this height.
	extension  *string // "jpg", "png"
	focalPoint *string // "top", "bottom", "center"
	quality    *int    // quality of the output image, only used for jpg
}

type batchConfig struct {
	inputPath       string
	outputPath      string
	outputExtension string
}

type imageToProcess struct {
	inputPath      string
	inputMimetype  string
	outputPath     string
	outputMimetype string
}

type resizerConfig struct {
	mimetype   string
	width      int
	height     int
	focalPoint string
	quality    int
}

func main() {
	input := inputFlags{
		inputPath:  flag.String("input", "", "Input path, can be a file or directory."),
		outputPath: flag.String("output", "", "The output directory path."),
		width:      flag.Int("width", 640, "Desired output width."),
		height:     flag.Int("height", -1, "Desired output height. If > -1, image is cropped to this height."),
		extension:  flag.String("extension", "jpg", "Output image extension. Options: jpg, png."),
		focalPoint: flag.String("crop", "center", "Crop position. Options: top, bottom, center."),
		quality:    flag.Int("quality", 95, "Quality of the output image. Only used for jpg."),
	}

	flag.Parse()

	// was the program called with no arguments?
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	failIf(verifyInuput(input))

	// get all files in the imagesPath directory
	imagesToProcess, err := getImagesToProcess(batchConfig{
		inputPath:       *input.inputPath,
		outputPath:      *input.outputPath,
		outputExtension: *input.extension})
	failIf(err)

	log.Printf("Found %d images to process\n", len(imagesToProcess))

	for _, imageToProcess := range imagesToProcess {
		log.Printf("Processing file:\n\tin: %s\n\tout: %s", imageToProcess.inputPath, imageToProcess.outputPath)

		original, err := os.Open(imageToProcess.inputPath)
		failIf(err)
		defer original.Close()

		resized, err := os.Create(imageToProcess.outputPath)
		failIf(err)
		defer resized.Close()

		decodedImage, err := getImage(original, imageToProcess.inputMimetype)
		failIf(err)

		resizeImage(decodedImage, resized, resizerConfig{
			mimetype:   imageToProcess.outputMimetype,
			width:      *input.width,
			height:     *input.height,
			focalPoint: *input.focalPoint,
			quality:    *input.quality})
	}
}

func resizeImage(src image.Image, resized io.Writer, config resizerConfig) error {
	inputBounds := src.Bounds()

	// Scale the image to the desired width while maintaining the aspect ratio
	aspectRatio := float64(inputBounds.Dy()) / float64(inputBounds.Dx())
	scaledHeight := int(math.Round(float64(config.width) * aspectRatio))
	scaledBounds := image.Rect(0, 0, config.width, scaledHeight)
	scaledImage := image.NewRGBA(scaledBounds)
	draw.BiLinear.Scale(scaledImage, scaledImage.Rect, src, inputBounds, draw.Over, nil)

	// Crop the scaled image to the desired height
	cropHeight := config.height
	if cropHeight > scaledHeight {
		cropHeight = scaledHeight // Ensure we don't crop beyond the scaled image
	}

	cropY := 0 // Default to top crop
	switch config.focalPoint {
	case "center":
		cropY = (scaledHeight - cropHeight) / 2
	case "bottom":
		cropY = scaledHeight - cropHeight
	}
	croppedBounds := image.Rect(0, cropY, config.width, cropY+cropHeight)
	croppedImage := scaledImage.SubImage(croppedBounds)

	// Step 3: Encode the cropped image to the desired format
	switch config.mimetype {
	case "image/png":
		return png.Encode(resized, croppedImage)
	default:
		return jpeg.Encode(resized, croppedImage, &jpeg.Options{Quality: config.quality})
	}
}

func getImage(original io.Reader, mimeType string) (image.Image, error) {
	var src image.Image
	var err error

	switch mimeType {
	case "image/jpeg":
		src, err = jpeg.Decode(original)
	case "image/png":
		src, err = png.Decode(original)
	case "image/webp":
		src, err = webp.Decode(original)
	default:
		src, err = nil, fmt.Errorf("unsupported image format: %s", mimeType)
	}

	if err != nil {
		return nil, err
	}

	return src, nil
}

func getImagesToProcess(config batchConfig) ([]imageToProcess, error) {
	inputStat, err := os.Stat(config.inputPath)
	if err != nil {
		return nil, err
	}

	var files []string
	if inputStat.IsDir() {
		dirFiles, err := os.ReadDir(config.inputPath)
		if err != nil {
			return nil, err
		}
		for _, file := range dirFiles {
			files = append(files, filepath.Join(config.inputPath, file.Name()))
		}
	} else {
		files = []string{config.inputPath}
	}

	var images []imageToProcess

	for _, file := range files {
		if isImage(file) {
			outputName := strings.Replace(filepath.Base(file), filepath.Ext(file), "."+config.outputExtension, 1)
			outputPath := filepath.Join(config.outputPath, outputName)

			images = append(images, imageToProcess{
				inputPath:      file,
				inputMimetype:  getMimetype(file),
				outputPath:     outputPath,
				outputMimetype: getMimetype(outputPath)})
		}
	}

	return images, nil
}

func getMimetype(path string) string {
	// get the file extension to determine the mimetype
	// log.Println("image extension: ", imageExtension)
	mimetype := "image/jpeg"
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		mimetype = "image/jpeg"
	case ".png":
		mimetype = "image/png"
	case ".webp":
		mimetype = "image/webp"
	}
	return mimetype
}

func isImage(path string) bool {
	// check if the file is an image by checking the file extension
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	default:
		return false
	}
}

func verifyInuput(f inputFlags) error {
	if *f.inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if _, err := os.Stat(*f.inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input path does not exist: %s", *f.inputPath)
	}
	if *f.outputPath == "" {
		return fmt.Errorf("output path is required")
	}
	if _, err := os.Stat(*f.outputPath); os.IsNotExist(err) {
		return fmt.Errorf("output path does not exist: %s", *f.outputPath)
	}
	if *f.width <= 0 {
		return fmt.Errorf("width must be greater than 0")
	}
	if *f.height < -1 || *f.height == 0 {
		return fmt.Errorf("height must be -1 or greater than 0")
	}
	if *f.extension != "jpg" && *f.extension != "png" {
		return fmt.Errorf("invalid extension, must be jpg or png")
	}
	if *f.focalPoint != "top" && *f.focalPoint != "bottom" && *f.focalPoint != "left" && *f.focalPoint != "right" && *f.focalPoint != "center" {
		return fmt.Errorf("invalid crop position, must be top, bottom, left, right or center")
	}
	if *f.quality < 0 || *f.quality > 100 {
		return fmt.Errorf("quality must be between 0 and 100")
	}

	return nil
}

func failIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
