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
            description: "The output file/directory to process.")
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

        var heightOption = new Option<int?>(
            aliases: ["--height", "-h"],
            description: "Desired output height.")
        {
            ArgumentHelpName = "If specified, image is cropped to this height post-scaling."
        };

        AddOption(heightOption);

        this.SetHandler((input, output, type, quality, width, height) =>
            ImageProcessor.Process(new(input, output, type, quality, width, height)),
            inputArgument, outputOption, imageOutputType, qualityOption, widthOption, heightOption);
    }
}

public sealed record CommandParams(
    string Input,
    string Output,
    ImageOutputType Type,
    int Quality,
    int Width,
    int? Height);

public enum ImageOutputType
{
    jpg,
    png
}
