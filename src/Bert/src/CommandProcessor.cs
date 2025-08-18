using Danom;

namespace Bert;

public static class CommandProcessor {
    public static void Process(string input, string output, ImageParams param) {
        IEnumerable<Result<ImageProcessingOutput, ImageProcessingError>> results;
        if (File.Exists(input)) {
            results = [ProcessFile(input, output, param)];
        }
        else if (Directory.Exists(input)) {
            results = ProcessDirectory(input, output, param);
        }
        else {
            results = [ImageProcessingResult.Error(input, "Input path does not exist")];
        }

        Print.ResultsTable(results);

        Console.WriteLine();
        Print.Info($"Output: {output}");
    }

    static IEnumerable<Result<ImageProcessingOutput, ImageProcessingError>> ProcessDirectory(
        string input,
        string output,
        ImageParams param) {
        if (!Directory.Exists(input)) {
            yield return ImageProcessingResult.Error(input, $"Input directory does not exist");
        }

        foreach (var filePath in Directory.EnumerateFiles(input)) {
            yield return ProcessFile(filePath, output, param with { FilePath = filePath });
        }
    }

    static Result<ImageProcessingOutput, ImageProcessingError> ProcessFile(string input, string output, ImageParams param) {
        var file = new FileInfo(input);

        if (!file.Exists) {
            return ImageProcessingResult.Error(file.FullName, "Input file does not exist");
        }

        if (!file.IsImage()) {
            return ImageProcessingResult.Error(file.FullName, "Input file is not a valid image");
        }

        string outputPath = output; // Assume output is a file path

        if (!Path.HasExtension(output)) {
            // output is a directory, create output file name based on input file
            var outputName = Path.ChangeExtension(Path.GetFileName(file.FullName), param.OutputType.ToString());
            outputPath = Path.Combine(output, outputName);

            if (!Directory.Exists(output)) {
                Directory.CreateDirectory(output);
            }
        }

        if (File.Exists(outputPath)) {
            File.Delete(outputPath);
        }

        using var inputStream = File.OpenRead(file.FullName);
        using var outputStream = File.OpenWrite(outputPath);

        return param switch {
            CropParams cropParams =>
                ImageProcessor.CropImage(cropParams, inputStream, outputStream),

            ScaleByWidthParams scaleByWidthParams =>
                ImageProcessor.ScaleImage(scaleByWidthParams, inputStream, outputStream),

            ScaleByHeightParams scaleByHeightParams =>
                ImageProcessor.ScaleImage(scaleByHeightParams, inputStream, outputStream),

            _ => ImageProcessingResult.Error(file.FullName, "Unsupported command mode.")
        };
    }
}
