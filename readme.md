<div align="center">

![bert](https://github.com/eastcitysoftware/bert/blob/assets/bert.png?raw=true)

Resize images in bulk.
</div>

---

Meet Bert, the no-nonsense tool for resizing and compressing images in bulk. Got a pile of oversized assets clogging up your workflow? Bert's got your back. Forget the endless clicks and convoluted methods—just tell Bert what you need, and he'll handle the rest. Resize, crop, compress—Bert does it all with minimal fuss and maximum efficiency.

Because life's too short for bloated images.

---

## Usage

```
Usage of .\bert.exe:
  -crop string
        Crop position. Options: top, bottom, center. (default "center")
  -extension string
        Output image extension. Options: jpg, png. (default "jpg")
  -height int
        Desired output height. If > -1, image is cropped to this height. (default -1)
  -input string
        Input path, can be a file or directory.
  -output string
        The output directory path.
  -quality int
        Quality of the output image. Only used for jpg. (default 95)
  -width int
        Desired output width. (default 640)
```

## Example

```bash
bert -input ./input -output ./output -width 800 -height 600 -quality 80 -extension jpg -crop center
```

This will resize all images in the input directory to 800x600, with a quality of 80, and save them as jpg files in the output directory. The images will be cropped to the center if they are not already 800x600.

## Example with a single file

```bash
bert -input ./input/image.png -output ./output -width 800 -height 600 -quality 80 -extension jpg -crop center
```

This will resize the image.png file to 800x600, with a quality of 80, and save it as a jpg file in the output directory. The image will be cropped to the center if it is not already 800x600.
