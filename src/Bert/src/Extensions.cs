namespace Bert;

public static class FileInfoExtensions {
    public static bool IsImage(this FileInfo file) =>
        file.Extension.Equals(".jpg", StringComparison.OrdinalIgnoreCase)
        || file.Extension.Equals(".jpeg", StringComparison.OrdinalIgnoreCase)
        || file.Extension.Equals(".png", StringComparison.OrdinalIgnoreCase)
        || file.Extension.Equals(".webp", StringComparison.OrdinalIgnoreCase);
}
