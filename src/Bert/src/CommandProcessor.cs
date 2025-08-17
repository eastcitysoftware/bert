using Danom;

namespace Bert;

public static class CommandProcessor {
    public static void Process(string input, string output, ImageParams param) {
        if (File.Exists(input)) {
            ProcessFile(input, output, param);
        }
        else if (Directory.Exists(input)) {
            ProcessDirectory(input, output, param);
        }
        else {
            Console.WriteLine($"Input path does not exist: {input}");
        }
    }

    static void ProcessDirectory(string input, string output, ImageParams param) {
        if (!Directory.Exists(input)) {
            Console.WriteLine($"Input directory does not exist: {input}");
            return;
        }

        foreach (var filePath in Directory.EnumerateFiles(input)) {
            ProcessFile(filePath, output, param);
        }
    }

    static void ProcessFile(string input, string output, ImageParams param) {
        var file = new FileInfo(input);

        if (!file.Exists) {
            Console.WriteLine($"Input file does not exist: {file.FullName}");
            return;
        }

        if (!file.IsImage()) {
            Console.WriteLine($"Input file is not a valid image: {file.FullName}");
            return;
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

        Console.WriteLine($"Processing {file.FullName}");

        if (File.Exists(outputPath)) {
            File.Delete(outputPath);
        }

        using var inputStream = File.OpenRead(file.FullName);
        using var outputStream = File.OpenWrite(outputPath);

        var outputResult = param switch {
            CropParams cropParams =>
                ImageProcessor.CropImage(cropParams, inputStream, outputStream),

            ScaleByWidthParams scaleByWidthParams =>
                ImageProcessor.ScaleImage(scaleByWidthParams, inputStream, outputStream),

            ScaleByHeightParams scaleByHeightParams =>
                ImageProcessor.ScaleImage(scaleByHeightParams, inputStream, outputStream),

            _ => Result.Error("Unsupported command mode.")
        };

        outputResult
            .Match(
                ok: _ => Console.WriteLine($"Image processed successfully: {outputPath}"),
                error: e => Console.WriteLine($"Error processing image: {e}"));
    }
}
