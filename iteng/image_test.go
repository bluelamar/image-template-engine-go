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
	"testing"
)

type min_test struct {
	a   int
	b   int
	exp int
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
	imagePath := "../resources/sun_and_moon_100x100.png"
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
	imagePath = "../resources/test_template.json"
	_, err = LoadImageFromFile(imagePath)
	if err != nil {
		t.Logf("Success: LoadImageFromFile(%s) returned error: %v", imagePath, err)
	} else {
		t.Errorf("LoadImageFromFile(%s) should have returned an error but didnt", imagePath)
	}
}

func Test_jpg_LoadImageFromFile(t *testing.T) {
	// Test loading a jpg image from a file
	imagePath := "../resources/sun_and_moon_100x100.jpg"
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
	imagePath := "../resources/sun_and_moon_100x100.tiff"
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
	imagePath := "../resources/sun_and_moon_100x100.png"
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
	imagePath := "../resources/sun_and_moon_100x100.jpg"
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
	imagePath := "../resources/sun_and_moon_100x100.tiff"
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
