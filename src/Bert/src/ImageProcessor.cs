using Danom;
using SkiaSharp;

namespace Bert;

public abstract record ImageParams(
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High);

public abstract record ScaleParams(
    int ScaleTo,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High) : ImageParams(OutputType, Quality);

public sealed record ScaleByWidthParams(
    int Width,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High,
    CropY CropY = CropY.Top,
    int? CropHeight = default) : ScaleParams(Width, OutputType, Quality);

public sealed record ScaleByHeightParams(
    int Height,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High,
    CropX CropX = CropX.Center,
    int? CropWidth = default) : ScaleParams(Height, OutputType, Quality);

public sealed record CropParams(
    Dimensions CropDimensions,
    CropX CropX,
    CropY CropY,
    ImageType OutputType = ImageType.jpg,
    ImageQuality Quality = ImageQuality.High) : ImageParams(OutputType, Quality);

public static class ImageProcessor {
    public static Result<Unit, ResultErrors> ScaleImage(
        ScaleParams param,
        Stream input,
        Stream output) {
        using var image = SKBitmap.Decode(input);

        if (image is null) {
            return Result.Error("Failed to decode input image.");
        }

        Dimensions? scaleDimensions = default;
        CropWindow? cropWindow = default;

        if (param is ScaleByWidthParams scaleByWidth) {
            // scale the image to the specified width
            var scaleByWithDimensions = Scaler.ScaleToWidth(scaleByWidth.ScaleTo, new(image.Width, image.Height));

            if (scaleByWidth.CropHeight is int cropHeight
                && cropHeight > 0
                && cropHeight < scaleByWithDimensions.Height) {
                cropWindow = Cropper.Crop(CropX.Left, scaleByWidth.CropY, scaleByWithDimensions with { Height = cropHeight }, scaleByWithDimensions);
            }

            scaleDimensions = scaleByWithDimensions;
        }
        else if (param is ScaleByHeightParams scaleByHeight) {
            // scale the image to the specified height
            var scaleByHeightDimensions = Scaler.ScaleToHeight(scaleByHeight.ScaleTo, new(image.Width, image.Height));

            if (scaleByHeight.CropWidth is int cropWidth
                && cropWidth > 0
                && cropWidth < scaleByHeightDimensions.Width) {
                cropWindow = Cropper.Crop(scaleByHeight.CropX, CropY.Top, scaleByHeightDimensions with { Width = cropWidth }, scaleByHeightDimensions);
            }

            scaleDimensions = scaleByHeightDimensions;
        }

        if (scaleDimensions is Dimensions dimensions
            && (dimensions.Width != image.Width || dimensions.Height != image.Height)) {
            // scale the image using the calculated dimensions
            Console.WriteLine($"\tScaling image from {image.Width}x{image.Height} to {dimensions.Width}x{dimensions.Height}");
            var scale = new SKImageInfo(dimensions.Width, dimensions.Height);
            using var scaledImage = image.Resize(scale, new SKSamplingOptions(SKFilterMode.Linear, SKMipmapMode.Linear));

            if (scaledImage is null) {
                return Result.Error("Failed to scale image.");
            }

            if (cropWindow is CropWindow crop) {
                // crop the scaled image
                Console.WriteLine($"\tCropping image using coordinates (L:{crop.Left}, T:{crop.Top}, R:{crop.Right}, B:{crop.Bottom})");
                var cropRect = new SKRectI(crop.Left, crop.Top, crop.Right, crop.Bottom);
                using var croppedImage = new SKBitmap();

                if (!scaledImage.ExtractSubset(croppedImage, cropRect)) {
                    return Result.Error("Failed to crop image.");
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

        return Result.Ok();
    }

    public static Result<Unit, ResultErrors> CropImage(
        CropParams param,
        Stream input,
        Stream output) {
        using var image = SKBitmap.Decode(input);

        if (image is null) {
            return Result.Error("Failed to decode input image.");
        }

        var cropWindow = Cropper.Crop(param.CropX, param.CropY, param.CropDimensions, new(image.Width, image.Height));
        Console.WriteLine($"\tCropping image using coordinates (L:{cropWindow.Left}, T:{cropWindow.Top}, R:{cropWindow.Right}, B:{cropWindow.Bottom})");
        using var croppedImage = new SKBitmap();

        if (!image.ExtractSubset(croppedImage, new(cropWindow.Left, cropWindow.Top, cropWindow.Right, cropWindow.Bottom))) {
            return Result.Error("Failed to extract subset from image.");
        }

        using var scaledImage = croppedImage.Resize(
            new SKImageInfo(param.CropDimensions.Width, param.CropDimensions.Height),
            new SKSamplingOptions(SKFilterMode.Linear, SKMipmapMode.Linear));

        scaledImage.Encode(output, GetImageFormat(param.OutputType), (int)param.Quality);

        return Result.Ok(Unit.Value);
    }

    private static SKEncodedImageFormat GetImageFormat(ImageType type) =>
        type switch {
            ImageType.png => SKEncodedImageFormat.Png,
            _ => SKEncodedImageFormat.Jpeg
        };
}
