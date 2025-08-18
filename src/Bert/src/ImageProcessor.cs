using Danom;
using SkiaSharp;

namespace Bert;

public abstract record ImageParams(
    string FilePath,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High);

public abstract record ScaleParams(
    string FilePath,
    int ScaleTo,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High) : ImageParams(FilePath, OutputType, Quality);

public sealed record ScaleByWidthParams(
    string FilePath,
    int Width,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High,
    CropY CropY = CropY.Top,
    int? CropHeight = default) : ScaleParams(FilePath, Width, OutputType, Quality);

public sealed record ScaleByHeightParams(
    string FilePath,
    int Height,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High,
    CropX CropX = CropX.Center,
    int? CropWidth = default) : ScaleParams(FilePath, Height, OutputType, Quality);

public sealed record CropParams(
    string FilePath,
    Dimensions CropDimensions,
    CropX CropX,
    CropY CropY,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High) : ImageParams(FilePath, OutputType, Quality);

public sealed record ImageProcessingOutput(
    string FilePath,
    Dimensions OriginalDimensions,
    Dimensions ProcessedDimensions,
    CropWindow? CropWindow = default);

public sealed record ImageProcessingError(
    string FilePath,
    string ErrorMessage);

public static class ImageProcessor {
    public static Result<ImageProcessingOutput, ImageProcessingError> ScaleImage(
        ScaleParams param,
        Stream input,
        Stream output) {
        using var image = SKBitmap.Decode(input);

        if (image is null) {
            return ImageProcessingResult.Error(param.FilePath, "Failed to decode input image.");
        }

        var imageDimensions = new Dimensions(image.Width, image.Height);
        Dimensions? scaleDimensions = default;
        CropWindow? cropWindow = default;

        if (param is ScaleByWidthParams scaleByWidth) {
            // scale the image to the specified width
            var scaleByWithDimensions = Scaler.ScaleToWidth(scaleByWidth.ScaleTo, imageDimensions);

            if (scaleByWidth.CropHeight is int cropHeight
                && cropHeight > 0
                && cropHeight < scaleByWithDimensions.Height) {
                cropWindow = Cropper.Crop(
                    cropX: CropX.Left,
                    cropY: scaleByWidth.CropY,
                    cropDimensions: scaleByWithDimensions with { Height = cropHeight },
                    dimensions: scaleByWithDimensions);
            }

            scaleDimensions = scaleByWithDimensions;
        }
        else if (param is ScaleByHeightParams scaleByHeight) {
            // scale the image to the specified height
            var scaleByHeightDimensions = Scaler.ScaleToHeight(scaleByHeight.ScaleTo, imageDimensions);

            if (scaleByHeight.CropWidth is int cropWidth
                && cropWidth > 0
                && cropWidth < scaleByHeightDimensions.Width) {
                cropWindow = Cropper.Crop(
                    cropX: scaleByHeight.CropX,
                    cropY: CropY.Top,
                    cropDimensions: scaleByHeightDimensions with { Width = cropWidth },
                    dimensions: scaleByHeightDimensions);
            }

            scaleDimensions = scaleByHeightDimensions;
        }

        if (scaleDimensions is Dimensions dimensions
            && dimensions != imageDimensions) {
            // scale the image using the calculated dimensions
            var scale = new SKImageInfo(dimensions.Width, dimensions.Height);
            using var scaledImage = image.Resize(scale, new SKSamplingOptions(SKFilterMode.Linear, SKMipmapMode.Linear));

            if (scaledImage is null) {
                return ImageProcessingResult.Error(param.FilePath, "Failed to scale image.");
            }

            if (cropWindow is CropWindow crop) {
                // crop the scaled image
                var cropRect = new SKRectI(crop.Left, crop.Top, crop.Right, crop.Bottom);
                using var croppedImage = new SKBitmap();

                if (!scaledImage.ExtractSubset(croppedImage, cropRect)) {
                    return ImageProcessingResult.Error(param.FilePath, "Failed to crop image.");
                }

                croppedImage.Encode(output, GetImageFormat(param.OutputType), (int)param.Quality);
            }
            else {
                // no crop window, just encode the scaled image
                scaledImage.Encode(output, GetImageFormat(param.OutputType), (int)param.Quality);
            }
        }
        else {
            // no scaling, just encode the original image
            image.Encode(output, GetImageFormat(param.OutputType), (int)param.Quality);
        }

        return ImageProcessingResult.Ok(param.FilePath, imageDimensions, scaleDimensions ?? cropWindow?.Dimensions ?? imageDimensions, cropWindow);
    }

    public static Result<ImageProcessingOutput, ImageProcessingError> CropImage(
        CropParams param,
        Stream input,
        Stream output) {
        using var image = SKBitmap.Decode(input);

        if (image is null) {
            return ImageProcessingResult.Error(param.FilePath, "Failed to decode input image.");
        }

        var imageDimensions = new Dimensions(image.Width, image.Height);
        var cropWindow = Cropper.Crop(param.CropX, param.CropY, param.CropDimensions, imageDimensions);
        using var croppedImage = new SKBitmap();

        if (!image.ExtractSubset(croppedImage, new(cropWindow.Left, cropWindow.Top, cropWindow.Right, cropWindow.Bottom))) {
            return ImageProcessingResult.Error(param.FilePath, "Failed to extract subset from image.");
        }

        using var scaledImage = croppedImage.Resize(
            new SKImageInfo(param.CropDimensions.Width, param.CropDimensions.Height),
            new SKSamplingOptions(SKFilterMode.Linear, SKMipmapMode.Linear));

        scaledImage.Encode(output, GetImageFormat(param.OutputType), (int)param.Quality);

        return ImageProcessingResult.Ok(param.FilePath, imageDimensions, cropWindow.Dimensions, cropWindow);
    }

    private static SKEncodedImageFormat GetImageFormat(ImageType type) =>
        type switch {
            ImageType.png => SKEncodedImageFormat.Png,
            _ => SKEncodedImageFormat.Jpeg
        };
}

public static class ImageProcessingResult {
    public static Result<ImageProcessingOutput, ImageProcessingError> Ok(
        string filePath,
        Dimensions originalDimensions,
        Dimensions processedDimensions,
        CropWindow? cropWindow) =>
        Result<ImageProcessingOutput, ImageProcessingError>.Ok(new(filePath, originalDimensions, processedDimensions, cropWindow));

    public static Result<ImageProcessingOutput, ImageProcessingError> Error(
        string filePath,
        string errorMessage) =>
        Result<ImageProcessingOutput, ImageProcessingError>.Error(new(filePath, errorMessage));
}
