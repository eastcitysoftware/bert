namespace Bert;

using Danom;
using SkiaSharp;

internal static class Scaler
{
    internal static Result<Unit, ResultErrors> ScaleImage(ScaleConfig config, Stream input, Stream output)
    {
        using var image = SKBitmap.Decode(input);

        if (image is null)
        {
            return Result.Error("Failed to decode input image.");
        }

        if (config.OutputWidth > 0)
        {
            var scale = GetScale(config.OutputWidth, image.Width, image.Height);

            Console.WriteLine($"\tScaling image from {image.Width}x{image.Height} to {scale.Width}x{scale.Height}");

            using var scaledImage = image.Resize(scale, new SKSamplingOptions(SKFilterMode.Linear, SKMipmapMode.Linear));

            if (scaledImage is null)
            {
                return Result.Error("Failed to scale image.");
            }

            if (config.OutputHeight is int outputHeight && outputHeight > 0 && outputHeight < scaledImage.Height)
            {
                Console.WriteLine($"\tCropping image height from {scale.Height} to {outputHeight}");
                var crop = new SKRectI(0, 0, scaledImage.Width, outputHeight);
                using var croppedImage = new SKBitmap();

                if(!scaledImage.ExtractSubset(croppedImage, crop))
                {
                    return Result.Error("Failed to crop image.");
                }

                croppedImage.Encode(output, GetImageFormat(config.OutputType), config.Quality);
            }
            else
            {
                scaledImage.Encode(output, GetImageFormat(config.OutputType), config.Quality);
            }
        }
        else
        {
            image.Encode(output, GetImageFormat(config.OutputType), config.Quality);
        }

        return Result.Ok();
    }

    private static SKImageInfo GetScale(int outputWidth, int width, int height)
    {
        var aspectRatio = (double)height / width;
        var newHeight = (int)Math.Round(outputWidth * aspectRatio);
        return new SKImageInfo(outputWidth, newHeight);
    }

    private static SKEncodedImageFormat GetImageFormat(ImageOutputType type) =>
        type switch
        {
            ImageOutputType.png => SKEncodedImageFormat.Png,
            _ => SKEncodedImageFormat.Jpeg
        };
}

internal readonly record struct ScaleConfig(
    int OutputWidth,
    int? OutputHeight,
    ImageOutputType OutputType,
    int Quality);
