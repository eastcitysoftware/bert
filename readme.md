# bert (birt - bulk image resizer tool)

bert is a nice little fellow who will help you resize images in bulk.

```
Usage of .\bert.exe:
  -crop string
        Crop position. Options: top, bottom, left, right, center. (default "center")
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

## example

```bash
bert -input ./input -output ./output -width 800 -height 600 -quality 80 -extension jpg -crop center
```

This will resize all images in the input directory to 800x600, with a quality of 80, and save them as jpg files in the output directory. The images will be cropped to the center if they are not already 800x600.

## example with a single file

```bash
bert -input ./input/image.png -output ./output -width 800 -height 600 -quality 80 -extension jpg -crop center
```

This will resize the image.png file to 800x600, with a quality of 80, and save it as a jpg file in the output directory. The image will be cropped to the center if it is not already 800x600.
