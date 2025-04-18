using System.CommandLine;
using Bert;

var command = new BertCommand("bert", "Resize images in bulk.");
command.Invoke(args);
