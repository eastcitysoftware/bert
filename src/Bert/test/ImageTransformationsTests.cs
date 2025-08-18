namespace Bert.Tests;

public sealed class ImageTransformationsTests {
    [Fact]
    public void Dimensions_StructuralEquality() {
        var dim1 = new Dimensions(100, 200);
        var dim2 = new Dimensions(100, 200);
        var dim3 = new Dimensions(200, 100);

        Assert.Equal(dim1, dim2);
        Assert.NotEqual(dim1, dim3);
    }

    [Fact]
    public void Dimensions_ToStringOverride() {
        var dimensions = new Dimensions(100, 200);
        Assert.Equal("100x200", dimensions.ToString());
    }
}

public sealed class ScalerTests {
    [Theory]
    [InlineData(100, 200, 50, 100, 25)]
    [InlineData(150, 300, 75, 150, 37)]
    [InlineData(200, 400, 100, 200, 50)]
    public void ScaleToWidth(int outputWidth, int inputWidth, int inputHeight, int expectedWidth, int expectedHeight) {
        var inputDimensions = new Dimensions(inputWidth, inputHeight);
        var scaledDimensions = Scaler.ScaleToWidth(outputWidth, inputDimensions);
        Assert.Equal(expectedWidth, scaledDimensions.Width);
        Assert.Equal(expectedHeight, scaledDimensions.Height);
    }

    [Theory]
    [InlineData(100, 200, 150, 133, 100)]
    [InlineData(150, 300, 225, 200, 150)]
    [InlineData(200, 400, 300, 266, 200)]
    public void ScaleToHeight(int outputHeight, int inputWidth, int inputHeight, int expectedWidth, int expectedHeight) {
        var inputDimensions = new Dimensions(inputWidth, inputHeight);
        var scaledDimensions = Scaler.ScaleToHeight(outputHeight, inputDimensions);
        Assert.Equal(expectedWidth, scaledDimensions.Width);
        Assert.Equal(expectedHeight, scaledDimensions.Height);
    }
}

public sealed class CropperTests {
    [Theory]
    [InlineData(CropX.Left, CropY.Top, 50, 50, 100, 100, 0, 0)]
    [InlineData(CropX.Left, CropY.Center, 50, 50, 100, 100, 0, 25)]
    [InlineData(CropX.Left, CropY.Bottom, 50, 50, 100, 100, 0, 50)]
    [InlineData(CropX.Center, CropY.Top, 50, 50, 100, 100, 25, 0)]
    [InlineData(CropX.Center, CropY.Center, 50, 50, 100, 100, 25, 25)]
    [InlineData(CropX.Center, CropY.Bottom, 50, 50, 100, 100, 25, 50)]
    [InlineData(CropX.Right, CropY.Top, 50, 50, 100, 100, 50, 0)]
    [InlineData(CropX.Right, CropY.Center, 50, 50, 100, 100, 50, 25)]
    [InlineData(CropX.Right, CropY.Bottom, 50, 50, 100, 100, 50, 50)]
    public void Crop(
        CropX cropX, CropY cropY,
        int cropWidth, int cropHeight,
        int imageWidth, int imageHeight,
        int expectedLeft, int expectedTop) {
        var cropDimensions = new Dimensions(cropWidth, cropHeight);
        var dimensions = new Dimensions(imageWidth, imageHeight);
        var cropWindow = Cropper.Crop(cropX, cropY, cropDimensions, dimensions);
        Assert.Equal(expectedLeft, cropWindow.Left);
        Assert.Equal(expectedTop, cropWindow.Top);
        Assert.Equal(expectedLeft + cropWidth, cropWindow.Right);
        Assert.Equal(expectedTop + cropHeight, cropWindow.Bottom);
    }
}
