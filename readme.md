<div align="center">

![bert](https://github.com/eastcitysoftware/bert/blob/assets/bert.png?raw=true)

[![build](https://github.com/eastcitysoftware/bert/actions/workflows/build.yml/badge.svg)](https://github.com/eastcitysoftware/bert/actions/workflows/build.yml)
![License](https://img.shields.io/github/license/eastcitysoftware/bert)

Resize images in bulk.
</div>

---

Meet Bert, the no-nonsense tool for resizing and compressing images in bulk. Got a pile of oversized assets clogging up your workflow? Bert's got your back. Forget the endless clicks and convoluted methods—just tell Bert what you need, and he'll handle the rest. Resize, crop, compress—Bert does it all with minimal fuss and maximum efficiency.

Because life's too short for bloated images.

---

## Usage

```
Usage of .\bert.exe:

Description:
  Resize images in bulk.

Usage:
  bert <input> [options]

Arguments:
  <input>  The input file/directory to process.

Options:
  -o, --output <output> (REQUIRED)                                            The output directory to process.
  -t, --type <jpg|png>                                                        The image output type to process. [default: jpg]
  -q, --quality <1-100 (default: 80)> (REQUIRED)                              The quality of the image to process. [default: 80]
  -w, --width <If not specified, no scaling will be done.>                    Desired output width.
  -h, --height <If specified, image is cropped to this height post-scaling.>  Desired output height.
  --version                                                                   Show version information
  -?, -h, --help                                                              Show help and usage information
```

## Example

```bash
bert ./input --output ./output --width 800 --height 600 --quality 80 --type jpg --crop top
```

This will resize all images in the input directory to 800x600, with a quality of 80, and save them as jpg files in the output directory. The images will be cropped to the top if they are not already 800x600.

## Example with a single file

```bash
bert ./input/image.png --output ./output --width 800 --height 600 --quality 80 --type jpg -crop center
```

This will resize the image.png file to 800x600, with a quality of 80, and save it as a jpg file in the output directory. The image will be cropped to the center if it is not already 800x600.
