using Xunit;

namespace Bert.Tests;

public sealed class FileInfoExtensionsTests {
    [Theory]
    [InlineData("test.jpg")]
    [InlineData("test.jpeg")]
    [InlineData("test.png")]
    [InlineData("test.webp")]
    public void IsImage_ValidImageFile_ReturnsTrue(string filePath) {
        var file = new FileInfo(filePath);
        Assert.True(file.IsImage());
    }

    [Theory]
    [InlineData("test.txt")]
    [InlineData("test.docx")]
    [InlineData("test.pdf")]
    [InlineData("test.mp3")]
    [InlineData("test")]
    public void IsImage_InvalidFile_ReturnsFalse(string filePath) {
        var file = new FileInfo(filePath);
        Assert.False(file.IsImage());
    }
}
