using SkiaSharp;

namespace Bert.Tests;

public sealed class ImageProcessorTests {
    private const string _filePath = "assets/logo.png";

    [Theory]
    [InlineData(100, null, ImageType.png, 100, 50)]
    [InlineData(100, 25, ImageType.png, 100, 25)]
    [InlineData(100, null, ImageType.jpg, 100, 50)]
    [InlineData(100, 25, ImageType.jpg, 100, 25)]
    public void ScaleImageByWidth(int width, int? cropHeight, ImageType outputType, int expectedWidth, int expectedHeight) {
        var config = new ScaleByWidthParams(_filePath, Width: width, CropHeight: cropHeight, OutputType: outputType);
        var outputFile = $"assets/scaleimagebywidth.{width}-{cropHeight ?? expectedHeight}.{outputType}";

        File.Delete(outputFile);
        using var inputStream = File.OpenRead(_filePath);
        using var outputStream = File.OpenWrite(outputFile);

        var result = ImageProcessor.ScaleImage(config, inputStream, outputStream);

        inputStream.Close();
        outputStream.Close();

        Assert.True(result.IsOk);
        Assert.True(File.Exists(outputFile));

        using var outputImage = SKBitmap.Decode(File.OpenRead(outputFile));
        Assert.Equal(expectedWidth, outputImage.Width);
        Assert.Equal(expectedHeight, outputImage.Height);
    }

    [Theory]
    [InlineData(100, null, ImageType.png, 200, 100)]
    [InlineData(100, 25, ImageType.png, 25, 100)]
    [InlineData(100, null, ImageType.jpg, 200, 100)]
    [InlineData(100, 25, ImageType.jpg, 25, 100)]
    public void ScaleImageByHeight(int height, int? cropWidth, ImageType outputType, int expectedWidth, int expectedHeight) {
        var config = new ScaleByHeightParams(_filePath, Height: height, CropWidth: cropWidth, OutputType: outputType);
        var outputFile = $"assets/scaleimagebyheight.{height}-{cropWidth ?? expectedWidth}.{outputType}";

        File.Delete(outputFile);
        using var inputStream = File.OpenRead(_filePath);
        using var outputStream = File.OpenWrite(outputFile);

        var result = ImageProcessor.ScaleImage(config, inputStream, outputStream);

        inputStream.Close();
        outputStream.Close();

        Assert.True(result.IsOk);
        Assert.True(File.Exists(outputFile));

        using var outputImage = SKBitmap.Decode(File.OpenRead(outputFile));
        Assert.Equal(expectedWidth, outputImage.Width);
        Assert.Equal(expectedHeight, outputImage.Height);
    }

    [Theory]
    [InlineData(100, 50, ImageType.png, 100, 50)]
    [InlineData(100, 25, ImageType.png, 100, 25)]
    [InlineData(100, 50, ImageType.jpg, 100, 50)]
    [InlineData(100, 25, ImageType.jpg, 100, 25)]
    public void CropImage (int width, int height, ImageType outputType, int expectedWidth, int expectedHeight) {
        var config = new CropParams(_filePath, new(width, height), CropX.Left, CropY.Top, outputType);
        var outputFile = $"assets/cropimage.{width}-{height}.{outputType}";

        File.Delete(outputFile);
        using var inputStream = File.OpenRead(_filePath);
        using var outputStream = File.OpenWrite(outputFile);

        var result = ImageProcessor.CropImage(config, inputStream, outputStream);

        inputStream.Close();
        outputStream.Close();

        Assert.True(result.IsOk);
        Assert.True(File.Exists(outputFile));

        using var outputImage = SKBitmap.Decode(File.OpenRead(outputFile));
        Assert.Equal(expectedWidth, outputImage.Width);
        Assert.Equal(expectedHeight, outputImage.Height);
    }
}
