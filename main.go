package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eastcitysoftware/bert/internal/crop"
	"github.com/eastcitysoftware/bert/internal/scale"
	"golang.org/x/image/webp"
)

const (
	jpegMimeType = "image/jpeg"
	pngMimeType  = "image/png"
	webpMimeType = "image/webp"
)

type imageToProcess struct {
	inputPath  string
	outputPath string
}

type resize struct {
	outputMimetype string
	width          int
	height         int
	focalPoint     string
	quality        int
}

func main() {
	inputPath := flag.String("input", "", "input path, can be a file or directory")
	outputPath := flag.String("output", "", "the output directory path")
	width := flag.Int("width", 640, "desired output width")
	height := flag.Int("height", -1, "desired output height. If > -1, image is cropped to this height")
	extension := flag.String("extension", "jpg", "output image extension. Options: jpg, png")
	focalPoint := flag.String("crop", "center", "crop position. Options: top, bottom, center")
	quality := flag.Int("quality", 95, "quality of the output image. Only used for jpg")

	flag.Parse()

	// was the program called with no arguments?
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	// was the program called with unnamed arguments?
	if len(os.Args) >= 3 {
		if *inputPath == "" {
			inputPath = &os.Args[1]
		}

		if *outputPath == "" {
			outputPath = &os.Args[2]
		}
	}

	// verify input
	if *inputPath == "" {
		log.Fatalf("input path is required")
	}
	if _, err := os.Stat(*inputPath); os.IsNotExist(err) {
		log.Fatalf("input path does not exist: %s", *inputPath)
	}
	if *outputPath == "" {
		log.Fatalf("output path is required")
	}
	if outputStat, err := os.Stat(*outputPath); !outputStat.IsDir() || os.IsNotExist(err) {
		log.Fatalf("output path is invalid, or does not exist: %s", *outputPath)
	}
	if *width <= 0 {
		log.Fatalf("width must be greater than 0")
	}
	if *height < -1 || *height == 0 {
		log.Fatalf("height must be -1 or greater than 0")
	}
	if *extension != "jpg" && *extension != "png" {
		log.Fatalf("invalid extension, must be jpg or png")
	}
	if *focalPoint != "top" && *focalPoint != "bottom" && *focalPoint != "left" && *focalPoint != "right" && *focalPoint != "center" {
		log.Fatalf("invalid crop position, must be top, bottom, left, right or center")
	}
	if *quality < 0 || *quality > 100 {
		log.Fatalf("quality must be between 0 and 100")
	}

	// get all files in the imagesPath directory
	imagesToProcess, err := getImages(*inputPath, *outputPath, *extension)
	if err != nil {
		log.Fatalf("failed to get images for processing: %v", imagesToProcess)
	}
	log.Printf("found %d images to process\n", len(imagesToProcess))

	resize := resize{
		outputMimetype: getMimetype(*extension),
		width:          *width,
		height:         *height,
		focalPoint:     *focalPoint,
		quality:        *quality}

	for _, imageToProcess := range imagesToProcess {
		log.Printf("processing file:\n\tin: %s\n\tout: %s", imageToProcess.inputPath, imageToProcess.outputPath)
		processImage(imageToProcess, resize)
	}
}

func processImage(src imageToProcess, resize resize) error {
	original, err := os.Open(src.inputPath)
	if err != nil {
		return fmt.Errorf("failed to open src file: %v", err)
	}
	defer original.Close()

	resized, err := os.Create(src.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create resize file: %v", err)
	}
	defer resized.Close()

	decodedImage, err := getImage(original, getMimetype(src.inputPath))
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	return resizeImage(decodedImage, resized, resize)
}

func resizeImage(src image.Image, dst io.Writer, config resize) error {
	// Scale the image to the desired width while maintaining the aspect ratio
	scaledImage := scale.ScaleImage(src, config.width)

	var cropImage image.Image
	if config.height > 0 {
		scaledBounds := scaledImage.Bounds()
		scaledWidth := scaledBounds.Dx()
		scaledHeight := scaledBounds.Dy()
		cropImage = crop.CropImage(
			scaledImage,
			scaledWidth,
			scaledHeight,
			min(config.height, scaledHeight),
			config.focalPoint)
	}

	switch {
	case config.outputMimetype == pngMimeType && cropImage != nil:
		return png.Encode(dst, cropImage)
	case config.outputMimetype == pngMimeType:
		return png.Encode(dst, scaledImage)
	case config.outputMimetype == jpegMimeType && cropImage != nil:
		return jpeg.Encode(dst, cropImage, &jpeg.Options{Quality: config.quality})
	default:
		return jpeg.Encode(dst, scaledImage, &jpeg.Options{Quality: config.quality})
	}
}

func getImage(original io.Reader, mimeType string) (image.Image, error) {
	var src image.Image
	var err error

	switch mimeType {
	case jpegMimeType:
		src, err = jpeg.Decode(original)
	case pngMimeType:
		src, err = png.Decode(original)
	case webpMimeType:
		src, err = webp.Decode(original)
	default:
		src, err = nil, fmt.Errorf("unsupported image format: %s", mimeType)
	}

	if err != nil {
		return nil, err
	}

	return src, nil
}

func getImages(inputPath string, outputPath string, outputExtension string) ([]imageToProcess, error) {
	inputStat, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}

	var files []string
	if inputStat.IsDir() {
		dirFiles, err := os.ReadDir(inputPath)
		if err != nil {
			return nil, err
		}
		for _, file := range dirFiles {
			files = append(files, filepath.Join(inputPath, file.Name()))
		}
	} else {
		files = []string{inputPath}
	}

	var images []imageToProcess

	for _, file := range files {
		if isImage(file) {
			outputName := strings.Replace(filepath.Base(file), filepath.Ext(file), "."+outputExtension, 1)
			outputPath := filepath.Join(outputPath, outputName)

			images = append(images, imageToProcess{
				inputPath:  file,
				outputPath: outputPath})
		}
	}

	return images, nil
}

func getMimetype(path string) string {
	// get the file extension to determine the mimetype
	// log.Println("image extension: ", imageExtension)
	mimetype := jpegMimeType
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		mimetype = jpegMimeType
	case ".png":
		mimetype = pngMimeType
	case ".webp":
		mimetype = webpMimeType
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
