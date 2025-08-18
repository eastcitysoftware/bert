namespace Bert;

public readonly record struct Dimensions(
    int Width,
    int Height) {
        public override string ToString() => $"{Width}x{Height}";
    }

public readonly record struct CropWindow(
    int Left,
    int Top,
    int Right,
    int Bottom,
    Dimensions Dimensions);

public enum CropX {
    Left,
    Center,
    Right
}

public enum CropY {
    Top,
    Center,
    Bottom
}

public enum ImageType {
    jpg,
    png
}

public enum ImageQuality {
    Full = 100,
    High = 90,
    Medium = 75,
    Low = 60
}

public static class Scaler {
    public static Dimensions ScaleToWidth(int width, Dimensions dimensions) {
        var aspectRatio = (double)dimensions.Height / dimensions.Width;
        var newHeight = (int)Math.Floor(width * aspectRatio);
        return new(width, newHeight);
    }

    public static Dimensions ScaleToHeight(int height, Dimensions dimensions) {
        var aspectRatio = (double)dimensions.Width / dimensions.Height;
        var newWidth = (int)Math.Floor(height * aspectRatio);
        return new(newWidth, height);
    }
}

public static class Cropper {
    public static CropWindow Crop(CropX cropX, CropY cropY, Dimensions cropDimensions, Dimensions dimensions) {
        int x = cropX switch {
            CropX.Left => 0,
            CropX.Center => (dimensions.Width - cropDimensions.Width) / 2,
            CropX.Right => dimensions.Width - cropDimensions.Width,
            _ => throw new ArgumentOutOfRangeException(nameof(cropX), cropX, null)
        };

        int y = cropY switch {
            CropY.Top => 0,
            CropY.Center => (dimensions.Height - cropDimensions.Height) / 2,
            CropY.Bottom => dimensions.Height - cropDimensions.Height,
            _ => throw new ArgumentOutOfRangeException(nameof(cropY), cropY, null)
        };

        return new CropWindow(x, y, cropDimensions.Width + x, cropDimensions.Height + y, cropDimensions);
    }
}
