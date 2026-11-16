using Danom;

namespace Bert;

public static class Print {
    public static void Error(string message) {
        WriteColoredMessage("ERROR", ConsoleColor.Red);
        Console.Write(" ");
        Console.Write(message);
        Console.WriteLine();
    }

    public static void Warning(string message) {
        WriteColoredMessage("WARNING", ConsoleColor.Yellow);
        Console.Write(" ");
        Console.Write(message);
        Console.WriteLine();
    }

    private static void WriteColoredMessage(string message, ConsoleColor color) {
        var originalColor = Console.BackgroundColor;
        Console.BackgroundColor = color;
        Console.Write(message);
        Console.BackgroundColor = originalColor;
    }

    public static void Info(string message)
    {
        Console.WriteLine(message);
    }

    public static void Block(string message) {
        int boxWidth = message.Length + 4;
        var border = new string('═', boxWidth - 2);
        Console.WriteLine();
        Console.WriteLine($"╔═{border}═╗");
        Console.WriteLine($"║ {message}{new string(' ', boxWidth - message.Length - 1)}║");
        Console.WriteLine($"╚═{border}═╝");
    }

    public static void ResultsTable(IEnumerable<Result<ImageProcessingOutput, ImageProcessingError>> results) {
        var outputList = new List<ImageProcessingOutput>();
        var errorList = new List<ImageProcessingError>();

        foreach (var result in results) {
            result.Match(
                ok: outputList.Add,
                error: errorList.Add);
        }

        var filePaths = outputList.Select(x => x.FilePath).Concat(errorList.Select(x => x.FilePath)).Distinct().ToList();
        int nameColWidth = Math.Max(15, filePaths.Max(x => x.Length));
        Console.WriteLine();
        Console.WriteLine(new string('=', nameColWidth + 40));
        Console.WriteLine($"{"File".PadRight(nameColWidth)} | {"Original",-15} | {"Result",-15}");
        Console.WriteLine(new string('-', nameColWidth + 40));

        foreach (var output in outputList) {
            Console.WriteLine($"{output.FilePath.PadRight(nameColWidth)} | " +
                             $"{output.OriginalDimensions,-15} | " +
                             $"{output.ProcessedDimensions,-15}");
        }

        foreach (var error in errorList) {
            Console.WriteLine($"{error.FilePath.PadRight(nameColWidth)} | " +
                             $"{string.Empty,-15} | " +
                             $"Error: {error.ErrorMessage}");
        }

        Console.WriteLine(new string('=', nameColWidth + 40));
    }
}
