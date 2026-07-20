# image-template-engine-go
Given a base image, create an updated image using a template of images and text.

## Summary

You start with a base image in a supported format, ie. png, jpg, webp, etc.
Then you can lay over it with smaller images and text to specified locations on the base image.
This will produce the final updated image (leaving the base image intact) with the required modifications.

### Template file

You specify a template file that defines **slots** in your image.
For each **slot** you specify an image or text, positioniong and other options.
See the example **slot** template in `test/example_template.json`.

### Input file

Each of the defined **slots** in the `template` file will contain an identifier field `id`.
Then you create an **input** file that will specify the content for each `id`
specified in the **template** file.
See the example **input** file in `test/example_input.json`.


### Supported Image Formats

* bmp
* gif
* jpg
* png
* tiff

**TODO** Need tests for bmp and gif

Resource image details:

* Image Size: 1024 X 1024 pixels
* RGB
* DPI : 72 pixels/inch


## Installation

```bash
# To get the latest released Go client:
go get github.com/bluelamar/image-template-engine-go@latest
```


## Usage


### Example

See https://github.com/bluelamar/image-template-engine-go/tree/main/examples

See the function **ImageDriver** in https://github.com/bluelamar/image-template-engine-go/tree/master/iteng/driver.go to see how to use the API.

The example main can be run with the script `bin/run_example.sh`.
