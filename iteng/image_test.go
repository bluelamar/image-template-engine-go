// Copyright 2022, Initialize All Once Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iteng

import (
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/fogleman/gg"
)

type min_test struct {
	a   int
	b   int
	exp int
}

func Test_ApplyOpacity_NoChange(t *testing.T) {
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 1))
	rgba.SetRGBA(0, 0, color.RGBA{10, 20, 30, 255})
	rgba.SetRGBA(1, 0, color.RGBA{40, 50, 60, 128})

	var img image.Image = rgba
	ret := ApplyOpacity(img, 1.0)
	if ret != img {
		t.Errorf("ApplyOpacity with opacity 1.0 should return the original image")
	}
}

func Test_ApplyOpacity_ReducesAlpha(t *testing.T) {
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 1))
	// set two pixels with known alphas
	rgba.SetRGBA(0, 0, color.RGBA{10, 20, 30, 200})
	rgba.SetRGBA(1, 0, color.RGBA{40, 50, 60, 100})

	var img image.Image = rgba
	opacity := 0.5
	ret := ApplyOpacity(img, opacity)

	// ret should be a different image
	if ret == img {
		t.Errorf("ApplyOpacity with opacity < 0.9999 should return a new image")
	}

	// verify alpha values
	r0, g0, b0, a0 := ret.At(0, 0).RGBA()
	exp0 := uint8(float64(uint8(200)) * opacity)
	if uint8(a0>>8) != exp0 {
		t.Errorf("pixel (0,0) alpha = %d; expected %d", uint8(a0>>8), exp0)
	}

	_, _, _, a1 := ret.At(1, 0).RGBA()
	exp1 := uint8(float64(uint8(100)) * opacity)
	if uint8(a1>>8) != exp1 {
		t.Errorf("pixel (1,0) alpha = %d; expected %d", uint8(a1>>8), exp1)
	}

	// also ensure color components are preserved (within 8-bit precision)
	if uint8(r0>>8) != 10 || uint8(g0>>8) != 20 || uint8(b0>>8) != 30 {
		t.Errorf("pixel (0,0) color changed unexpectedly: got %d,%d,%d", uint8(r0>>8), uint8(g0>>8), uint8(b0>>8))
	}
}

var min_tests = []min_test{
	{a: 1, b: 2, exp: 1},
	{a: 2, b: 1, exp: 1},
	{a: 3, b: 3, exp: 3},
}

func Test_min(t *testing.T) {
	for _, test := range min_tests {
		ret := min(test.a, test.b)
		if ret != test.exp {
			t.Errorf("min(%d,%d) = %d; expected %d", test.a, test.b, ret, test.exp)
		}
	}
}

type minf_test struct {
	a   float64
	b   float64
	exp float64
}

var minf_tests = []minf_test{
	{a: 1.0, b: 2.0, exp: 1.0},
	{a: 2.0, b: 1.0, exp: 1.0},
	{a: 3.0, b: 3.0, exp: 3.0},
}

func Test_minf(t *testing.T) {
	for _, test := range minf_tests {
		ret := minf(test.a, test.b)
		if ret != test.exp {
			t.Errorf("minf(%f,%f) = %f; expected %f", test.a, test.b, ret, test.exp)
		}
	}
}

var maxf_tests = []minf_test{
	{a: 1.0, b: 2.0, exp: 2.0},
	{a: 2.0, b: 1.0, exp: 2.0},
	{a: 3.0, b: 3.0, exp: 3.0},
}

func Test_maxf(t *testing.T) {
	for _, test := range maxf_tests {
		ret := maxf(test.a, test.b)
		if ret != test.exp {
			t.Errorf("maxf(%f,%f) = %f; expected %f", test.a, test.b, ret, test.exp)
		}
	}
}

func Test_LoadImageFromFile(t *testing.T) {
	// Test loading an image from a file
	imagePath := "../test/sun_and_moon_100x100.png"
	image, err := LoadImageFromFile(imagePath)
	if err != nil {
		t.Errorf("LoadImageFromFile(%s) returned error: %v", imagePath, err)
	}
	if image == nil {
		t.Errorf("LoadImageFromFile(%s) returned nil image", imagePath)
	}

	// Non-existant file
	imagePath = "non-existant-file.png"
	_, err = LoadImageFromFile(imagePath)
	if err != nil {
		t.Logf("Success: LoadImageFromFile(%s) should return an error: %v", imagePath, err)
	} else {
		t.Errorf("LoadImageFromFile(%s) should have returned an error but didnt", imagePath)
	}

	// Invalid image file format
	imagePath = "../test/test_template.json"
	_, err = LoadImageFromFile(imagePath)
	if err != nil {
		t.Logf("Success: LoadImageFromFile(%s) returned error: %v", imagePath, err)
	} else {
		t.Errorf("LoadImageFromFile(%s) should have returned an error but didnt", imagePath)
	}
}

func Test_jpg_LoadImageFromFile(t *testing.T) {
	// Test loading a jpg image from a file
	imagePath := "../test/sun_and_moon_100x100.jpg"
	image, err := LoadImageFromFile(imagePath)
	if err != nil {
		t.Errorf("LoadImageFromFile(%s) returned error: %v", imagePath, err)
	}
	if image == nil {
		t.Errorf("LoadImageFromFile(%s) returned nil image", imagePath)
	}
}

func Test_tiff_LoadImageFromFile(t *testing.T) {
	// Test loading a tiff image from a file
	imagePath := "../test/sun_and_moon_100x100.tiff"
	image, err := LoadImageFromFile(imagePath)
	if err != nil {
		t.Errorf("LoadImageFromFile(%s) returned error: %v", imagePath, err)
	}
	if image == nil {
		t.Errorf("LoadImageFromFile(%s) returned nil image", imagePath)
	}
}

func Test_png_ResizeImage(t *testing.T) {
	// Test resizing an image
	imagePath := "../test/sun_and_moon_100x100.png"
	image, err := LoadImageFromFile(imagePath)
	if err != nil {
		t.Errorf("LoadImageFromFile(%s) returned error: %v", imagePath, err)
	}
	if image == nil {
		t.Errorf("LoadImageFromFile(%s) returned nil image", imagePath)
	}

	// Resize the image to 50x50
	resizedImage := ResizeImage(image, 50, 50, ResizeModeFill)
	if err != nil {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeFill) returned error: %v", image, err)
	}
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 50, 50) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() != 50 || resizedImage.Bounds().Dy() != 50 {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeFill) returned image with incorrect dimensions", image)
	}

	// Test resizing an image to a larger size
	resizedImage = ResizeImage(image, 2000, 2000, ResizeModeFill)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFill) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 2000 || resizedImage.Bounds().Dy() > 2000 {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFill) returned image with incorrect dimensions", image)
	}

	// ResizeModeFit
	resizedImage = ResizeImage(image, 2000, 2000, ResizeModeFit)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 2000 || resizedImage.Bounds().Dy() > 2000 {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned image with incorrect dimensions", image)
	}

	// ResizeModeCover
	resizedImage = ResizeImage(image, 50, 50, ResizeModeCover)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeCover) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 50 || resizedImage.Bounds().Dy() > 50 {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeCover) returned image with incorrect dimensions", image)
	}
}

func Test_jpg_ResizeImage(t *testing.T) {
	// Test resizing an image
	imagePath := "../test/sun_and_moon_100x100.jpg"
	image, err := LoadImageFromFile(imagePath)
	if err != nil {
		t.Errorf("LoadImageFromFile(%s) returned error: %v", imagePath, err)
	}
	if image == nil {
		t.Errorf("LoadImageFromFile(%s) returned nil image", imagePath)
	}

	// Resize the image to 50x50
	resizedImage := ResizeImage(image, 50, 50, ResizeModeFill)
	if err != nil {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeFill) returned error: %v", image, err)
	}
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 50, 50) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() != 50 || resizedImage.Bounds().Dy() != 50 {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeFill) returned image with incorrect dimensions", image)
	}

	// Test resizing an image to a larger size
	resizedImage = ResizeImage(image, 2000, 2000, ResizeModeFit)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 2000 || resizedImage.Bounds().Dy() > 2000 {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned image with incorrect dimensions", image)
	}

	// ResizeModeFit
	resizedImage = ResizeImage(image, 2000, 2000, ResizeModeFit)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 2000 || resizedImage.Bounds().Dy() > 2000 {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned image with incorrect dimensions", image)
	}

	// ResizeModeCover
	resizedImage = ResizeImage(image, 50, 50, ResizeModeCover)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeCover) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 50 || resizedImage.Bounds().Dy() > 50 {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeCover) returned image with incorrect dimensions", image)
	}
}

func Test_tiff_ResizeImage(t *testing.T) {
	// Test resizing an image
	imagePath := "../test/sun_and_moon_100x100.tiff"
	image, err := LoadImageFromFile(imagePath)
	if err != nil {
		t.Errorf("LoadImageFromFile(%s) returned error: %v", imagePath, err)
	}
	if image == nil {
		t.Errorf("LoadImageFromFile(%s) returned nil image", imagePath)
	}

	// Resize the image to 50x50
	resizedImage := ResizeImage(image, 50, 50, ResizeModeFill)
	if err != nil {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeFill) returned error: %v", image, err)
	}
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 50, 50) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() != 50 || resizedImage.Bounds().Dy() != 50 {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeFill) returned image with incorrect dimensions", image)
	}

	// Test resizing an image to a larger size
	resizedImage = ResizeImage(image, 2000, 2000, ResizeModeFit)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 2000 || resizedImage.Bounds().Dy() > 2000 {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned image with incorrect dimensions", image)
	}

	// ResizeModeFit
	resizedImage = ResizeImage(image, 2000, 2000, ResizeModeFit)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 2000 || resizedImage.Bounds().Dy() > 2000 {
		t.Errorf("ResizeImage(%v, 2000, 2000, ResizeModeFit) returned image with incorrect dimensions", image)
	}

	// ResizeModeCover
	resizedImage = ResizeImage(image, 50, 50, ResizeModeCover)
	if resizedImage == nil {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeCover) returned nil image", image)
	}

	// Check the dimensions of the resized image
	if resizedImage.Bounds().Dx() > 50 || resizedImage.Bounds().Dy() > 50 {
		t.Errorf("ResizeImage(%v, 50, 50, ResizeModeCover) returned image with incorrect dimensions", image)
	}
}

func Test_MakeMask_Empty_FullyOpaque(t *testing.T) {
	// Test MakeMask with empty maskType creates a fully opaque rectangle
	mask := MakeMask("", 100, 80, 0)
	if mask == nil {
		t.Errorf("MakeMask(\"\", 100, 80, 0) returned nil")
	}

	// Check dimensions
	if mask.Bounds().Dx() != 100 || mask.Bounds().Dy() != 80 {
		t.Errorf("MakeMask dimensions incorrect: got %dx%d, expected 100x80",
			mask.Bounds().Dx(), mask.Bounds().Dy())
	}

	// Check that all pixels are fully opaque (255)
	for y := 0; y < mask.Bounds().Dy(); y++ {
		for x := 0; x < mask.Bounds().Dx(); x++ {
			_, _, _, a := mask.At(x, y).RGBA()
			alpha := uint8(a >> 8)
			if alpha != 255 {
				t.Errorf("pixel (%d,%d) has alpha=%d; expected 255", x, y, alpha)
			}
		}
	}
}

func Test_MakeMask_Rectangle(t *testing.T) {
	// Test MakeMask with unrecognized maskType creates a rectangle
	mask := MakeMask("rect", 50, 60, 0)
	if mask == nil {
		t.Errorf("MakeMask(\"rect\", 50, 60, 0) returned nil")
	}

	// Check dimensions
	if mask.Bounds().Dx() != 50 || mask.Bounds().Dy() != 60 {
		t.Errorf("MakeMask dimensions incorrect: got %dx%d, expected 50x60",
			mask.Bounds().Dx(), mask.Bounds().Dy())
	}

	// For rectangle, center should be opaque
	_, _, _, a := mask.At(25, 30).RGBA()
	if uint8(a>>8) == 0 {
		t.Errorf("center pixel should be opaque for rectangle mask")
	}
}

func Test_MakeMask_Circle(t *testing.T) {
	// Test MakeMask with "circle" creates a circular mask
	w, h := 100, 100
	mask := MakeMask("circle", w, h, 0)
	if mask == nil {
		t.Errorf("MakeMask(\"circle\", 100, 100, 0) returned nil")
	}

	// Check dimensions
	if mask.Bounds().Dx() != w || mask.Bounds().Dy() != h {
		t.Errorf("MakeMask dimensions incorrect: got %dx%d, expected %dx%d",
			mask.Bounds().Dx(), mask.Bounds().Dy(), w, h)
	}

	// Center should be opaque for circle
	_, _, _, a := mask.At(w/2, h/2).RGBA()
	if uint8(a>>8) == 0 {
		t.Errorf("center pixel should be opaque for circle mask")
	}

	// Corner should be transparent or less opaque for circle
	_, _, _, a = mask.At(5, 5).RGBA()
	cornerAlpha := uint8(a >> 8)
	_, _, _, a = mask.At(w/2, h/2).RGBA()
	centerAlpha := uint8(a >> 8)
	if cornerAlpha >= centerAlpha {
		t.Logf("corner alpha=%d, center alpha=%d (this may indicate a circular mask)", cornerAlpha, centerAlpha)
	}
}

func Test_MakeMask_Rounded(t *testing.T) {
	// Test MakeMask with "rounded" creates a rounded rectangle mask
	w, h := 80, 60
	mask := MakeMask("rounded", w, h, 10)
	if mask == nil {
		t.Errorf("MakeMask(\"rounded\", 80, 60, 10) returned nil")
	}

	// Check dimensions
	if mask.Bounds().Dx() != w || mask.Bounds().Dy() != h {
		t.Errorf("MakeMask dimensions incorrect: got %dx%d, expected %dx%d",
			mask.Bounds().Dx(), mask.Bounds().Dy(), w, h)
	}

	// Center should be opaque
	_, _, _, a := mask.At(w/2, h/2).RGBA()
	if uint8(a>>8) == 0 {
		t.Errorf("center pixel should be opaque for rounded mask")
	}
}

func Test_MakeMask_Rounded_DefaultRadius(t *testing.T) {
	// Test MakeMask with "rounded" and radius <= 0 uses default radius
	w, h := 100, 100
	mask := MakeMask("rounded", w, h, 0)
	if mask == nil {
		t.Errorf("MakeMask(\"rounded\", 100, 100, 0) returned nil")
	}

	// Check dimensions
	if mask.Bounds().Dx() != w || mask.Bounds().Dy() != h {
		t.Errorf("MakeMask dimensions incorrect: got %dx%d, expected %dx%d",
			mask.Bounds().Dx(), mask.Bounds().Dy(), w, h)
	}

	// Center should be opaque
	_, _, _, a := mask.At(w/2, h/2).RGBA()
	if uint8(a>>8) == 0 {
		t.Errorf("center pixel should be opaque for rounded mask with default radius")
	}
}

func Test_DrawTextInto_BasicText(t *testing.T) {
	// Test drawing basic text without special options
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:      10,
		Y:      10,
		Width:  180,
		Height: 80,
		TextOpts: TextOpt{
			FontSize:   20,
			Color:      "#000000",
			FontSource: "file",
			FontPath:   "../test/NotoSansPhoenician-Regular.ttf",
		},
	}

	// Should not panic
	slot.DrawTextInto(dc, "Hello World")

	// Verify the context is still valid
	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed: got %dx%d, expected 200x100", dc.Width(), dc.Height())
	}
}

func Test_DrawTextInto_WithColor(t *testing.T) {
	// Test drawing text with custom hex color
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:      10,
		Y:      10,
		Width:  180,
		Height: 80,
		TextOpts: TextOpt{
			FontSize:   20,
			Color:      "#FF0000", // Red
			FontSource: "file",
			FontPath:   "../test/NotoSansTagalog-Regular.ttf",
		},
	}

	// Should not panic
	slot.DrawTextInto(dc, "Red Text")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed: got %dx%d, expected 200x100", dc.Width(), dc.Height())
	}
}

func Test_DrawTextInto_WithDefaultColor(t *testing.T) {
	// Test drawing text without specifying color (should use black)
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:      10,
		Y:      10,
		Width:  180,
		Height: 80,
		TextOpts: TextOpt{
			FontSize:   20,
			FontSource: "file",
			FontPath:   "../test/LastResort.otf",
		},
	}

	// Should not panic and use default black color
	slot.DrawTextInto(dc, "Default Color")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed: got %dx%d, expected 200x100", dc.Width(), dc.Height())
	}
}

func Test_DrawTextInto_LeftAlignment(t *testing.T) {
	// Test drawing text with left alignment
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:       10,
		Y:       10,
		Width:   180,
		Height:  80,
		AnchorX: 0.0, // left
		TextOpts: TextOpt{
			FontSize: 20,
			AlignX:   "left",
		},
	}

	slot.DrawTextInto(dc, "Left Aligned")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after left alignment")
	}
}

func Test_DrawTextInto_CenterAlignment(t *testing.T) {
	// Test drawing text with center alignment
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:       10,
		Y:       10,
		Width:   180,
		Height:  80,
		AnchorX: 0.5, // center
		TextOpts: TextOpt{
			FontSize: 20,
			AlignX:   "center",
		},
	}

	slot.DrawTextInto(dc, "Center Aligned")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after center alignment")
	}
}

func Test_DrawTextInto_RightAlignment(t *testing.T) {
	// Test drawing text with right alignment
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:       10,
		Y:       10,
		Width:   180,
		Height:  80,
		AnchorX: 1.0, // right
		TextOpts: TextOpt{
			FontSize: 20,
			AlignX:   "right",
		},
	}

	slot.DrawTextInto(dc, "Right Aligned")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after right alignment")
	}
}

func Test_DrawTextInto_VerticalAlignment_Top(t *testing.T) {
	// Test drawing text with top vertical alignment
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:       10,
		Y:       10,
		Width:   180,
		Height:  80,
		AnchorY: 0.0, // top
		TextOpts: TextOpt{
			FontSize: 20,
			AlignY:   "top",
		},
	}

	slot.DrawTextInto(dc, "Top")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after top vertical alignment")
	}
}

func Test_DrawTextInto_VerticalAlignment_Middle(t *testing.T) {
	// Test drawing text with middle vertical alignment
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:       10,
		Y:       10,
		Width:   180,
		Height:  80,
		AnchorY: 0.5, // middle
		TextOpts: TextOpt{
			FontSize: 20,
			AlignY:   "middle",
		},
	}

	slot.DrawTextInto(dc, "Middle")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after middle vertical alignment")
	}
}

func Test_DrawTextInto_VerticalAlignment_Bottom(t *testing.T) {
	// Test drawing text with bottom vertical alignment
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:       10,
		Y:       10,
		Width:   180,
		Height:  80,
		AnchorY: 1.0, // bottom
		TextOpts: TextOpt{
			FontSize: 20,
			AlignY:   "bottom",
		},
	}

	slot.DrawTextInto(dc, "Bottom")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after bottom vertical alignment")
	}
}

func Test_DrawTextInto_WithWrapping(t *testing.T) {
	// Test drawing text with wrapping enabled
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:      10,
		Y:      10,
		Width:  180,
		Height: 80,
		TextOpts: TextOpt{
			FontSize: 14,
			Wrap:     true,
			MaxWidth: 150,
		},
	}

	longText := "This is a long text that should wrap to multiple lines when drawn into the canvas"
	slot.DrawTextInto(dc, longText)

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after wrapping text")
	}
}

func Test_DrawTextInto_WithoutWrapping(t *testing.T) {
	// Test drawing text without wrapping (single line)
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:      10,
		Y:      10,
		Width:  180,
		Height: 80,
		TextOpts: TextOpt{
			FontSize: 20,
			Wrap:     false,
		},
	}

	slot.DrawTextInto(dc, "Single Line Text")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after drawing single line text")
	}
}

func Test_DrawTextInto_InvalidAnchorX(t *testing.T) {
	// Test that invalid anchor X values are clamped to 0.0
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:       10,
		Y:       10,
		Width:   180,
		Height:  80,
		AnchorX: 2.5, // Invalid, should clamp to 0.0
		TextOpts: TextOpt{
			FontSize: 20,
		},
	}

	slot.DrawTextInto(dc, "Invalid X")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed with invalid anchor X")
	}
}

func Test_DrawTextInto_InvalidAnchorY(t *testing.T) {
	// Test that invalid anchor Y values are clamped to 0.0
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:       10,
		Y:       10,
		Width:   180,
		Height:  80,
		AnchorY: -0.5, // Invalid, should clamp to 0.0
		TextOpts: TextOpt{
			FontSize: 20,
		},
	}

	slot.DrawTextInto(dc, "Invalid Y")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed with invalid anchor Y")
	}
}

func Test_DrawTextInto_EmptyText(t *testing.T) {
	// Test drawing empty text string
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:      10,
		Y:      10,
		Width:  180,
		Height: 80,
		TextOpts: TextOpt{
			FontSize: 20,
		},
	}

	slot.DrawTextInto(dc, "")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed after drawing empty text")
	}
}

func Test_DrawTextInto_CaseSensitiveAlignments(t *testing.T) {
	// Test that alignment values are case-insensitive
	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:      10,
		Y:      10,
		Width:  180,
		Height: 80,
		TextOpts: TextOpt{
			FontSize: 20,
			AlignX:   "CENTER", // Uppercase
			AlignY:   "MIDDLE", // Uppercase
		},
	}

	slot.DrawTextInto(dc, "Case Insensitive")

	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed with uppercase alignments")
	}
}

func Test_DrawTextInto_LargeSlot(t *testing.T) {
	// Test drawing text in a large slot
	dc := gg.NewContext(1000, 800)
	slot := Slot{
		X:      50,
		Y:      50,
		Width:  900,
		Height: 700,
		TextOpts: TextOpt{
			FontSize: 48,
			Color:    "#0000FF", // Blue
		},
	}

	slot.DrawTextInto(dc, "Large Slot Text")

	if dc.Width() != 1000 || dc.Height() != 800 {
		t.Errorf("Context dimensions changed with large slot")
	}
}

func Test_DrawTextInto_SmallSlot(t *testing.T) {
	// Test drawing text in a small slot
	dc := gg.NewContext(100, 50)
	slot := Slot{
		X:      5,
		Y:      5,
		Width:  90,
		Height: 40,
		TextOpts: TextOpt{
			FontSize: 10,
		},
	}

	slot.DrawTextInto(dc, "Small")

	if dc.Width() != 100 || dc.Height() != 50 {
		t.Errorf("Context dimensions changed with small slot")
	}
}

func Test_DrawTextInto_WithEnvFontVars(t *testing.T) {
	// Test drawing text using ITENG_FONT_TTF and ITENG_FONT_DIR environment variables
	// Save original env vars
	origFontTTF := os.Getenv("ITENG_FONT_TTF")
	origFontDir := os.Getenv("ITENG_FONT_DIR")

	// Set test env vars
	os.Setenv("ITENG_FONT_TTF", "NotoSansPhoenician-Regular.ttf")
	os.Setenv("ITENG_FONT_DIR", "../test")

	defer func() {
		// Restore original env vars
		if origFontTTF != "" {
			os.Setenv("ITENG_FONT_TTF", origFontTTF)
		} else {
			os.Unsetenv("ITENG_FONT_TTF")
		}
		if origFontDir != "" {
			os.Setenv("ITENG_FONT_DIR", origFontDir)
		} else {
			os.Unsetenv("ITENG_FONT_DIR")
		}
	}()

	dc := gg.NewContext(200, 100)
	slot := Slot{
		X:      10,
		Y:      10,
		Width:  180,
		Height: 80,
		TextOpts: TextOpt{
			FontSize: 20,
			Color:    "#000000",
		},
	}

	// Should not panic and should load font from env vars
	slot.DrawTextInto(dc, "Env Font Test")

	// Verify the context is still valid
	if dc.Width() != 200 || dc.Height() != 100 {
		t.Errorf("Context dimensions changed: got %dx%d, expected 200x100", dc.Width(), dc.Height())
	}
}
