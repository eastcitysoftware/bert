namespace Bert;

using System.CommandLine;

public sealed class BertCommand : Command
{
    public BertCommand(string name, string description) : base(name, description)
    {
        var inputArgument = new Argument<string>(
            name: "input",
            description: "The input file/directory to process.");

        AddArgument(inputArgument);

        var outputOption = new Option<string>(
            aliases: ["--output", "-o"],
            description: "The output directory to process.")
        {
            IsRequired = true
        };

        AddOption(outputOption);

        var imageOutputType =
            new Option<ImageOutputType>(
                aliases: ["--type", "-t"],
                getDefaultValue: () => ImageOutputType.jpg,
                description: "The image output type to process.")
            {
                ArgumentHelpName = "jpg|png"
            };

        imageOutputType.AddCompletions("jpg", "png");

        AddOption(imageOutputType);

        var qualityOption =
            new Option<int>(
                aliases: ["--quality", "-q"],
                getDefaultValue: () => 80,
                description: "The quality of the image to process.")
            {
                IsRequired = true,
                ArgumentHelpName = "1-100 (default: 80)"
            };

        qualityOption.AddValidator(x =>
        {
            var value = x.GetValueOrDefault<int>();
            if (value < 1 || value > 100)
            {
                x.ErrorMessage = $"Invalid quality '{value}'. Valid values are between 1 and 100.";
            }
        });

        AddOption(qualityOption);

        var widthOption = new Option<int>(
            aliases: ["--width", "-w"],
            description: "Desired output width.")
        {
            ArgumentHelpName = "If not specified, no scaling will be done."
        };

        AddOption(widthOption);

        var heightOption = new Option<int>(
            aliases: ["--height", "-h"],
            description: "Desired output height.")
        {
            ArgumentHelpName = "If specified, image is cropped to this height post-scaling."
        };

        AddOption(heightOption);
        this.SetHandler((input, output, type, quality, width, height) =>
        {
            if (!Directory.Exists(output))
            {
                Directory.CreateDirectory(output);
            }

            // Process image(s)
            var scaleConfig = new ScaleConfig(
                OutputWidth: width,
                OutputHeight: height,
                OutputType: type,
                Quality: quality);

            if (File.Exists(input))
            {
                // Process single file
                var inputStream = File.OpenRead(input);
                var outputName = Path.ChangeExtension(Path.GetFileName(input), type.ToString());
                var outputPath = Path.Combine(output, outputName);

                Console.WriteLine($"Processing {input}");

                if (File.Exists(outputPath))
                {
                    File.Delete(outputPath);
                }

                using var outputStream = File.OpenWrite(outputPath);

                Scaler
                    .ScaleImage(scaleConfig, inputStream, outputStream)
                    .Match(
                        ok: _ => Console.WriteLine($"\tImage processed successfully: {outputPath}"),
                        error: x => Console.WriteLine($"\tError processing image: {x}"));
            }
            else if (Directory.Exists(input))
            {
                // Process directory
                foreach (var file in Directory.EnumerateFiles(input))
                {
                    if (!file.EndsWith(".jpg", StringComparison.OrdinalIgnoreCase)
                        && !file.EndsWith(".jpeg", StringComparison.OrdinalIgnoreCase)
                        && !file.EndsWith(".png", StringComparison.OrdinalIgnoreCase)
                        && !file.EndsWith(".webp", StringComparison.OrdinalIgnoreCase))
                    {
                        continue;
                    }

                    var inputStream = File.OpenRead(file);
                    var outputName = Path.ChangeExtension(Path.GetFileName(file), type.ToString());
                    var outputPath = Path.Combine(output, outputName);

                    Console.WriteLine($"Processing {file}");

                    if (File.Exists(outputPath))
                    {
                        File.Delete(outputPath);
                    }

                    using var outputStream = File.OpenWrite(outputPath);

                    Scaler
                        .ScaleImage(scaleConfig, inputStream, outputStream)
                        .Match(
                            ok: _ => Console.WriteLine($"\tOutput: {outputPath}"),
                            error: x => Console.WriteLine($"\tError: {x}"));
                }
            }
        }, inputArgument, outputOption, imageOutputType, qualityOption, widthOption, heightOption);
    }
}

public enum ImageOutputType
{
    jpg,
    png
}
