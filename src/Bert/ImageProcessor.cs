namespace Bert;

internal static class ImageProcessor
{
    internal static void Process(CommandParams param)
    {
        // Process image(s)
        var scaleConfig = new ScaleConfig(
            OutputWidth: param.Width,
            OutputHeight: param.Height,
            OutputType: param.Type,
            Quality: param.Quality);

        if (File.Exists(param.Input))
        {
            // Process single file
            string outputPath = param.Output;

            if (!Path.HasExtension(param.Output))
            {
                // output is a directory
                var outputName = Path.ChangeExtension(Path.GetFileName(param.Input), param.Type.ToString());
                outputPath = Path.Combine(param.Output, outputName);

                if (!Directory.Exists(param.Output))
                {
                    Directory.CreateDirectory(param.Output);
                }
            }

            Console.WriteLine($"Processing {param.Input}");

            var inputStream = File.OpenRead(param.Input);

            if (File.Exists(outputPath))
            {
                File.Delete(outputPath);
            }

            using var outputStream = File.OpenWrite(outputPath);

            Scaler.ScaleImage(scaleConfig, inputStream, outputStream)
                .Match(
                    ok: _ => Console.WriteLine($"\tImage processed successfully: {outputPath}"),
                    error: x => Console.WriteLine($"\tError processing image: {x}"));
        }
        else if (Directory.Exists(param.Input))
        {
            // Process directory
            if (!Directory.Exists(param.Output))
            {
                Directory.CreateDirectory(param.Output);
            }

            foreach (var file in Directory.EnumerateFiles(param.Input))
            {
                if (!IsImageFile(file))
                {
                    continue;
                }

                var inputStream = File.OpenRead(file);
                var outputName = Path.ChangeExtension(Path.GetFileName(file), param.Type.ToString());
                var outputPath = Path.Combine(param.Output, outputName);

                Console.WriteLine($"Processing {file}");

                if (File.Exists(outputPath))
                {
                    File.Delete(outputPath);
                }

                using var outputStream = File.OpenWrite(outputPath);

                Scaler.ScaleImage(scaleConfig, inputStream, outputStream)
                    .Match(
                        ok: _ => Console.WriteLine($"\tOutput: {outputPath}"),
                        error: x => Console.WriteLine($"\tError: {x}"));
            }
        }
    }

    private static bool IsImageFile(string file) =>
        file.EndsWith(".jpg", StringComparison.OrdinalIgnoreCase)
           || file.EndsWith(".jpeg", StringComparison.OrdinalIgnoreCase)
           || file.EndsWith(".png", StringComparison.OrdinalIgnoreCase)
           || file.EndsWith(".webp", StringComparison.OrdinalIgnoreCase);
}
