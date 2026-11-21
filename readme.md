<div align="center">

![bert](https://github.com/eastcitysoftware/bert/blob/assets/bert.png?raw=true)

[![build](https://github.com/eastcitysoftware/bert/actions/workflows/build.yml/badge.svg)](https://github.com/eastcitysoftware/bert/actions/workflows/build.yml)
![License](https://img.shields.io/github/license/eastcitysoftware/bert)

Resize images in bulk.
</div>

---

## Usage

```text
Usage of .\bert.exe:
  -crop-x string
        crop position. Options: left, right, center (default "center")
  -crop-y string
        crop position. Options: top, bottom, center (default "center")
  -height int
        desired output height. If > -1, image is cropped to this height (default -1)
  -input string
        input path, can be a file or directory
  -mode string
        resize mode. Options: scale, crop (default "scale")
  -output string
        the output directory path
  -quality string
        quality of the output image (jpg only). Options: full, high, medium and low (default "high")
  -type string
        image output type. Options: jpg, png (default "jpg")
  -width int
        desired output width (default 640)
```

The input and output paths can also be supplied as positional arguments. Place
all flags before positional paths:

```bash
bert -mode crop -width 543 -height 250 -quality low -type png ./input ./output
```

The equivalent command using named path flags is:

```bash
bert -input ./input -output ./output -mode crop -width 543 -height 250 -quality low -type png
```

The output directory is created automatically when it does not exist.

## Resize example

```bash
bert -input ./input -output ./output -mode scale -width 800 -type jpg -quality high
```

This scales every supported image in the input directory to a width of 800
pixels while preserving its aspect ratio. If `-height` is supplied in scale
mode, the image is scaled to that height instead.

## Crop example

```bash
bert -input ./input -output ./output -mode crop -width 800 -height 600 -crop-x center -crop-y center -type jpg -quality high
```

This crops every supported image to 800x600, using the center of the source
image. Use `-crop-x left|center|right` and
`-crop-y top|center|bottom` to select the crop position.

## Supported formats

The CLI processes files with these extensions:

```text
.jpg, .jpeg, .png, .webp
```

The output type is selected with `-type` and is either `jpg` or `png`.
JPEG quality is selected with `-quality`:

```text
full, high, medium, low
```

Quality affects JPEG output only.

## Find a bug?

There's an [issue](https://github.com/eastcitysoftware/bert/issues) for that.

## License

Licensed under [MIT](https://github.com/eastcitysoftware/bert/blob/main/LICENSE).