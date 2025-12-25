# image-template-engine-go
Given a base image, create an updated image using a template of images and text.

## Summary

You start with a base image in a supported format, ie. png, jpg, webp, etc.
You can layover it smaller images and text to specified locations on the base image.
This will produce the final updated image (leaving the base image untact) with the required modifications.

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

See https://github.com/bluelamar/image-template-engine-go/tree/master/examples

See the function **ImageDriver** in https://github.com/bluelamar/image-template-engine-go/tree/master/iteng/driver.go to see how to use the API.


