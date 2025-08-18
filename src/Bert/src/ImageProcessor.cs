using Danom;
using SixLabors.ImageSharp;
using SixLabors.ImageSharp.Formats.Png;
using SixLabors.ImageSharp.Formats.Jpeg;
using SixLabors.ImageSharp.Processing;
using SixLabors.ImageSharp.Formats;

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
        using var image = Image.Load(input);

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
            var scaleByHeightDimensions = Scaler.ScaleToHeight(scaleByHeight.ScaleTo, new(image.Width, image.Height));

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
            && (dimensions.Width != image.Width || dimensions.Height != image.Height)) {
            // scale the image using the calculated dimensions
            Console.WriteLine($"\tScaling image from {image.Width}x{image.Height} to {dimensions.Width}x{dimensions.Height}");

            image.Mutate(x =>
                x.Resize(width: dimensions.Width, height: dimensions.Height, sampler: KnownResamplers.Bicubic));

            if (cropWindow is CropWindow crop) {
                // crop the scaled image
                Console.WriteLine($"\tCropping image using coordinates (L:{crop.Left}, T:{crop.Top}, R:{crop.Right}, B:{crop.Bottom})");
                image.Mutate(x =>
                    x.Crop(new(
                        x: crop.Left,
                        y: crop.Top,
                        width: crop.Right - crop.Left,
                        height: crop.Bottom - crop.Top)));
            }
        }

        SaveImage(image, output, param.OutputType, param.Quality);

        return Result.Ok();
    }

    public static Result<Unit, ResultErrors> CropImage(
        CropParams param,
        Stream input,
        Stream output) {
        using var image = Image.Load(input);

        if (image is null) {
            return Result.Error("Failed to decode input image.");
        }

        var cropWindow = Cropper.Crop(param.CropX, param.CropY, param.CropDimensions, new(image.Width, image.Height));
        Console.WriteLine($"\tCropping image using coordinates (L:{cropWindow.Left}, T:{cropWindow.Top}, R:{cropWindow.Right}, B:{cropWindow.Bottom})");

        image.Mutate(x =>
            x.Crop(new(
                x: cropWindow.Left,
                y: cropWindow.Top,
                width: cropWindow.Right - cropWindow.Left,
                height: cropWindow.Bottom - cropWindow.Top)));

        SaveImage(image, output, param.OutputType, param.Quality);

        return Result.Ok(Unit.Value);
    }

    public static void SaveImage(
        Image image,
        Stream output,
        ImageType type,
        ImageQuality quality = ImageQuality.High) {
        switch (type) {
            case ImageType.png:
                image.SaveAsPng(output, GetImageEncoder(type, quality) as PngEncoder);
                break;
            default:
                image.SaveAsJpeg(output, GetImageEncoder(type, quality) as JpegEncoder);
                break;
        }
    }

    private static ImageEncoder GetImageEncoder(ImageType type, ImageQuality quality) =>
        type switch {
            ImageType.png => new PngEncoder { CompressionLevel = PngCompressionLevel.BestCompression },
            _ => new JpegEncoder { Quality = (int)quality }
        };
}
