namespace Bert;

using Danom;
using SkiaSharp;

internal readonly record struct ScaleConfig(
    int OutputWidth,
    int? OutputHeight,
    ImageOutputType OutputType,
    int Quality);

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

            scaledImage.Encode(output, config.OutputType switch
            {
                ImageOutputType.png => SKEncodedImageFormat.Png,
                _ => SKEncodedImageFormat.Jpeg
            }, config.Quality);
        }
        else
        {
            image.Encode(output, config.OutputType switch
            {
                ImageOutputType.png => SKEncodedImageFormat.Png,
                _ => SKEncodedImageFormat.Jpeg
            }, config.Quality);
        }

        return Result.Ok();
    }

    internal static SKImageInfo GetScale(int outputWidth, int width, int height)
    {
        var aspectRatio = (double)height / width;
        var newHeight = (int)Math.Round(outputWidth * aspectRatio);
        return new SKImageInfo(outputWidth, newHeight);
    }
}
