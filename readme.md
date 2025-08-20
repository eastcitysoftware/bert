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

```shell
Usage of .\bert.exe:

Description:
  bert, resize images in bulk

Usage:
  bert [command] [options]

Options:
  -?, -h, --help  Show help and usage information
  --version       Show version information

Commands:
  scale <input>  Scale images to the specified width/height.
  crop <input>   Crop images to the specified width/height.
```

### Scale Command

```shell
Description:
  Scale images to the specified width/height.

Usage:
  bert scale <input> [options]

Arguments:
  <input>  Input file/directory to process.

Options:
  output, -o (REQUIRED)               Output file/directory to process.
  type, -t <jpg|png>                  Output type. [default: jpg]
  quality, -q <Full|High|Medium|Low>  Output quality. [default: High]
  width, -w                           Desired output width in pixels.
  height, -h                          Desired output height in pixels.
  crop-x <Left|Center|Right>          Crop X position. [default: Center]
  crop-y <Top|Center|Bottom>          Crop Y position. [default: Top]
  -?, -h, --help                      Show help and usage information
```

### Crop Command

```shell
Description:
  Crop images to the specified width/height.

Usage:
  bert crop <input> [options]

Arguments:
  <input>  Input file/directory to process.

Options:
  output, -o (REQUIRED)       Output file/directory to process.
  type, -t <jpg|png>          Output type. [default: jpg]
  quality, -q <Full|High|Medium|Low>  Output quality. [default: High]
  width, -w                   Desired output width in pixels.
  height, -h                  Desired output height in pixels.
  crop-x <Left|Center|Right>  Crop X position. [default: Center]
  crop-y <Top|Center|Bottom>  Crop Y position. [default: Top]
  -?, -h, --help              Show help and usage information
```

## Examples

### Scaling a directory of images

```shell
bert scale ./input --output ./output --width 800 --height 600
```

This will resize all images in the input directory to 800x600, with a quality of 80, and save them as jpg files in the output directory. The images will be cropped to the top if they are not already 800x600.

## Scaling a single file

```shell
bert scale ./input/image.png --output ./output --width 800 --height 600
```

This will resize the image.png file to 800x600, with a quality of 80, and save it as a jpg file in the output directory. The image will be cropped to the center if it is not already 800x600.

## Scaling a single file with rename

```shell
bert scale ./input/image.png --output ./output/image-800x600.png --width 800 --height 600
```

This will resize the image.png file to 800x600, with a quality of 80, and save it as a jpg file in the output directory. The image will be cropped to the center if it is not already 800x600.
