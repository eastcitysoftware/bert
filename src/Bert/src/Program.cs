using System.CommandLine;
using Bert;

var command =
    new RootCommand(
        description: "bert, resize images in bulk");

var inputArgument =
    new Argument<string>(
        name: "input") {
        Description = "Input file/directory to process."
    };

var outputOption =
    new Option<string>("output", ["--output", "-o"]) {
        Description = "Output file/directory to process.",
        Required = true
    };

var outputTypeOption =
    new Option<ImageType>("type", ["--type", "-t"]) {
        DefaultValueFactory = _ => ImageType.jpg,
        Description = "Output type.",
        HelpName = "jpg|png",
        CompletionSources = { "jpg", "png" }
    };

var qualityOption =
    new Option<ImageQuality>("quality", ["--quality", "-q"]) {
        DefaultValueFactory = _ => ImageQuality.High,
        Description = "Output quality.",
        HelpName = "Full|High|Medium|Low",
        CompletionSources = { "Full", "High", "Medium", "Low" }
    };

var widthOption =
    new Option<int?>("width", ["--width", "-w"]) {
        Description = "Desired output width in pixels.",
    };
widthOption.Validators.Add(v => {
    if (v.GetValue(widthOption) is int w && w <= 0) {
        v.AddError("Width must be greater than zero.");
    }
});

var heightOption =
    new Option<int?>("height", ["--height", "-h"]) {
        Description = "Desired output height in pixels."
    };
heightOption.Validators.Add(v => {
    if (v.GetValue(heightOption) is int h && h <= 0) {
        v.AddError("Height must be greater than zero.");
    }
});

var cropXOption =
    new Option<CropX>("crop-x", ["--crop-x"]) {
        Description = "Crop X position.",
        DefaultValueFactory = _ => CropX.Center,
        HelpName = "Left|Center|Right",
        CompletionSources = { "Left", "Center", "Right" }
    };

var cropYOption =
    new Option<CropY>("crop-y", ["--crop-y"]) {
        Description = "Crop Y position.",
        DefaultValueFactory = _ => CropY.Top,
        HelpName = "Top|Center|Bottom",
        CompletionSources = { "Top", "Center", "Bottom" }
    };

var scaleCommand =
    new Command(
        name: "scale",
        description: "Scale images to the specified width/height.");

scaleCommand.Arguments.Add(inputArgument);
scaleCommand.Options.Add(outputOption);
scaleCommand.Options.Add(outputTypeOption);
scaleCommand.Options.Add(qualityOption);
scaleCommand.Options.Add(widthOption);
scaleCommand.Options.Add(heightOption);
scaleCommand.Options.Add(cropXOption);
scaleCommand.Options.Add(cropYOption);

scaleCommand.Validators.Add(commandResult => {
    if (commandResult.GetValue(inputArgument) is null
        || string.IsNullOrWhiteSpace(commandResult.GetValue(inputArgument))) {
        commandResult.AddError("Input argument is required.");
    }

    if (commandResult.GetValue(outputOption) is null
        || string.IsNullOrWhiteSpace(commandResult.GetValue(outputOption))) {
        commandResult.AddError("You must specify --output");
    }

    if (commandResult.GetValue(widthOption) is null
        && commandResult.GetValue(heightOption) is null) {
        commandResult.AddError("You must specify at least one of --width or --height.");
    }
});

var cropCommand =
    new Command(
        name: "crop",
        description: "Crop images to the specified width/height.");

cropCommand.Arguments.Add(inputArgument);
cropCommand.Options.Add(outputOption);
cropCommand.Options.Add(outputTypeOption);
cropCommand.Options.Add(widthOption);
cropCommand.Options.Add(heightOption);
cropCommand.Options.Add(cropXOption);
cropCommand.Options.Add(cropYOption);

cropCommand.Validators.Add(commandResult => {
    if (commandResult.GetValue(inputArgument) is null
        || string.IsNullOrWhiteSpace(commandResult.GetValue(inputArgument))) {
        commandResult.AddError("Input argument is required.");
    }

    if (commandResult.GetValue(outputOption) is null
        || string.IsNullOrWhiteSpace(commandResult.GetValue(outputOption))) {
        commandResult.AddError("Output option is required.");
    }

    if (commandResult.GetValue(widthOption) is not int w || w <= 0) {
        commandResult.AddError("You must specify a valid --width greater than zero.");
    }

    if (commandResult.GetValue(heightOption) is not int h || h <= 0) {
        commandResult.AddError("You must specify a valid --height greater than zero.");
    }
});

scaleCommand.SetAction(ExecuteCommand);
cropCommand.SetAction(ExecuteCommand);
command.SetAction(ExecuteCommand);
command.Subcommands.Add(scaleCommand);
command.Subcommands.Add(cropCommand);
command.Parse(args).Invoke();

int ExecuteCommand(ParseResult parseResult) {
    var input = parseResult.GetValue(inputArgument) ?? string.Empty;
    var output = parseResult.GetValue(outputOption) ?? string.Empty;
    var outputType = parseResult.GetValue(outputTypeOption);
    var quality = parseResult.GetValue(qualityOption);
    var width = parseResult.GetValue(widthOption);
    var height = parseResult.GetValue(heightOption);
    var cropX = parseResult.GetValue(cropXOption);
    var cropY = parseResult.GetValue(cropYOption);

    switch (parseResult.CommandResult.Command.Name, width, height) {
        case ("crop", int w, int h):
            Console.WriteLine($"Performing 'crop': {w}px x {h}px / {outputType} / {quality} / crop-x:{cropX} / crop-y:{cropY}");
            CommandProcessor.Process(input, output,
                new CropParams(
                    CropDimensions: new(w, h),
                    CropX: cropX,
                    CropY: cropY,
                    OutputType: outputType,
                    Quality: quality));

            break;

        case ("scale", int w, null):
            Console.WriteLine($"Performing 'scale': {w}px width / {outputType} / {quality} / crop-y:{cropY}");
            CommandProcessor.Process(input, output,
                new ScaleByWidthParams(
                    Width: w,
                    OutputType: outputType,
                    Quality: quality,
                    CropY: cropY,
                    CropHeight: null));
            break;

        case ("scale", null, int h):
            Console.WriteLine($"Performing 'scale': {h}px height / {outputType} / {quality} / crop-x:{cropX}");
            CommandProcessor.Process(input, output,
                new ScaleByHeightParams(
                    Height: h,
                    OutputType: outputType,
                    Quality: quality,
                    CropX: cropX,
                    CropWidth: null));
            break;

        case (_, int w, int h):
            Console.WriteLine($"Performing 'scale': {w}px x {h}px / {outputType} / {quality} / crop-x:{cropX} / crop-y:{cropY}");
            CommandProcessor.Process(input, output,
                new ScaleByWidthParams(
                    Width: w,
                    OutputType: outputType,
                    Quality: quality,
                    CropY: cropY,
                    CropHeight: h));
            break;

        default:
            Console.WriteLine("Invalid command or parameters. Use 'bert --help' for usage information.");
            break;
    }

    return 0;
}
