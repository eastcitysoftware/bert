package bert

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strings"
	"time"

	"github.com/eastcitysoftware/danom/check"
	"github.com/eastcitysoftware/danom/result"
	"golang.org/x/image/webp"
)

type TransformOptions struct {
	ImageSource   ImageSource
	OutputHandler OutputHandler
	Transform
}

type TransformResult struct {
	Transform
	Duration time.Duration
	Results  []ImageTransformResult
}

type ImageTransformResult struct {
	Name       string
	Dimensions Dimensions
}

func RunTransform(opts TransformOptions) result.Result[TransformResult, error] {
	validator := newTransformValidator()
	validation := validator.Validate(opts.Transform)

	if e, isErr := validation.ErrorValue(); isErr {
		return result.Error[TransformResult, error](e)
	}

	transform, isOk := validation.Value()

	if !isOk {
		return result.Error[TransformResult](fmt.Errorf("Invalid transform provided"))
	}

	startTime := time.Now()

	var results []ImageTransformResult
	var transformErrors []error
	for {
		nextImg := opts.ImageSource.Next()

		if nextImg.IsNone() {
			break
		}
		var transformer imageTransformer

		switch {
		case opts.Mode == TransformModeScale && opts.Width > 0 && opts.Height <= 0:
			transformer = transformToWidth(transform)
		case opts.Mode == TransformModeScale && opts.Height > 0:
			transformer = transformToHeight(transform)
		case opts.Mode == TransformModeCrop:
			transformer = transformToDimensions(transform)

		default:
			return result.Error[TransformResult](fmt.Errorf("unsupported transform"))
		}

		unprImg, _ := nextImg.Value()
		imgTransform := result.Map(unprImg.decode(), func(i image.Image) imageTransform {
			return opts.newImageTransform(unprImg.name, i)
		})

		result.Bind(imgTransform, transformer.transform).Do(
			func(i processedImage) {
				opts.OutputHandler.Handle(i).IterError(func(err error) {
					transformErrors = append(transformErrors, err)
				})

				results = append(results, ImageTransformResult{
					Name:       i.name,
					Dimensions: i.dimensions,
				})
			},
			func(e error) {
				transformErrors = append(transformErrors, e)
			})
	}

	if len(transformErrors) > 0 {
		return result.Error[TransformResult](errors.Join(transformErrors...))
	}

	return result.Ok[error](TransformResult{
		Transform: transform,
		Duration:  time.Since(startTime),
		Results:   results})
}

func newTransformValidator() check.Validator[Transform] {
	return check.NewValidator(
		check.NewRule("Mode", func(x Transform) string { return x.Mode },
			strAnyOf(TransformModeScale, TransformModeCrop)),
		check.NewRule("OutputType", func(x Transform) string { return x.OutputType },
			strAnyOf(JpgOutputType, PngOutputType)),
		check.NewRule("Quality", func(x Transform) string { return x.Quality },
			strAnyOf(FullQuality, HighQuality, MediumQuality, LowQuality)),
		validateUsing("Dimensions", func(x Transform) Dimensions { return x.Dimensions },
			newDimensionsValidator()),
		validateUsing("Crop", func(x Transform) Crop { return x.Crop },
			newCropValidator()),
	)
}

func newDimensionsValidator() check.Validator[Dimensions] {
	return check.NewValidator(
		check.NewRule("Dimensions", func(x Dimensions) Dimensions { return x },
			check.NewFact(func(x Dimensions) bool { return x.Width > 0 || x.Height > 0 },
				check.NewMessage("Width or height must be greater than zero"))),
	)
}

func newCropValidator() check.Validator[Crop] {
	return check.NewValidator(
		check.NewRule("CropX", func(x Crop) string { return x.CropX },
			strAnyOf(CropXCenter, CropXLeft, CropXRight)),
		check.NewRule("CropY", func(x Crop) string { return x.CropY },
			strAnyOf(CropYBottom, CropYCenter, CropYTop)))
}

func validateUsing[T, U any](name string, selector func(T) U, validator check.Validator[U]) check.Rule[T] {
	return func(input T) check.Error {
		result := validator(selector(input))
		nested, ok := result.ErrorValue()
		if !ok {
			return check.Error{Name: name}
		}

		e := check.Error{Name: name}
		for _, child := range nested {
			for _, message := range child.Errors {
				e.Errors = append(e.Errors, child.Name+" - "+message)
			}
		}

		return e
	}
}

func strAnyOf(vals ...string) check.Fact[string] {
	facts := []check.Fact[string]{}
	for _, v := range vals {
		facts = append(facts, check.StrEqualsFold(v))
	}
	return check.AnyOf(facts...)
}

func (t Transform) newImageTransform(name string, i image.Image) imageTransform {
	newName := strings.Replace(name, filepath.Ext(name), "", 1)
	switch t.OutputType {
	case PngOutputType:
		newName += PngExtension
	default:
		newName += JpgExtension
	}

	return imageTransform{
		image:     i,
		name:      newName,
		transform: t}
}

func (unprocessed unprocessedImage) decode() result.Result[image.Image, error] {
	rd := unprocessed.source()
	defer rd.Close()

	var src image.Image
	var err error
	ext := filepath.Ext(unprocessed.name)

	switch ext {
	case JpgExtension, JpegExtension:
		src, err = jpeg.Decode(rd)
	case PngExtension:
		src, err = png.Decode(rd)
	case WebpExtension:
		src, err = webp.Decode(rd)
	default:
		src, err = nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return result.Error[image.Image](err)
	}

	return result.Ok[error](src)
}
