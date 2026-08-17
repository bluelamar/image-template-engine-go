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
	"fmt"
	"image"
	"image/draw"
	"log"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

/*
ImageDriver is the main entry point for generating an image based on a template and inputs. It performs the following steps:
1. Parses the template JSON file to get the base image, output options, and slot definitions.
2. Parses the inputs JSON file to get the content for each slot.
3. Loads the base image and creates a canvas for drawing.
4. Iterates over each slot defined in the template:
  - If it's a text slot, it draws the text using the specified options.
  - If it's an image slot, it loads the input image, resizes it according to the slot's mode, applies opacity and masks if needed,
    and draws it onto the canvas at the correct position.

5. Saves the final generated image to the specified output path in the desired format.

This function abstracts away all the details of how images are loaded, resized, masked, and how text is drawn,
providing a simple interface for users to generate images based on templates and inputs.

Returns an error if any step fails, along with the final output path of the generated image.
*/
func ImageDriver(templatePath string, inputsPath string, outputPath string) (error, string) {
	tmpl, err := ParseTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err), ""
	}

	inputs, err := ParseInputsFromPath(inputsPath)
	if err != nil {
		return fmt.Errorf("parsing inputs: %v", err), ""
	}

	err, outputPath = GenerateImage(tmpl, inputs, outputPath)
	if err != nil {
		return fmt.Errorf("generating image: %v", err), ""
	}

	return nil, outputPath
}

/*
 * ImageDriverFromInputBytes is similar to ImageDriver, but it takes the inputs as a byte slice instead of a file path.
 * This allows for more flexibility in how the inputs are provided, such as reading from an HTTP request body or other sources.
 *
 * The function performs the same steps as ImageDriver, but it uses ParseInputsFromBytes to parse the inputs directly from the provided byte slice.
 *
 * Returns an error if any step fails, along with the final output path of the generated image.
 */
func ImageDriverFromInputBytes(templatePath string, inputsBytes []byte, outputPath string) (error, string) {
	tmpl, err := ParseTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err), ""
	}

	inputs, err := ParseInputsFromBytes(inputsBytes)
	if err != nil {
		return fmt.Errorf("parsing inputs: %v", err), ""
	}

	err, outputPath = GenerateImage(tmpl, inputs, outputPath)
	if err != nil {
		return fmt.Errorf("generating image: %v", err), ""
	}

	return nil, outputPath
}

/*
 * GenerateImage takes a parsed Template, Inputs, and an output path, and generates the final image according to the template and inputs.
 * It handles loading the base image, creating a canvas, processing each slot (text or image), applying resizing, opacity,
 * and masks as needed, and finally saving the output image.
 *
 * Note that if the output format is not specified in the template, it will be inferred from the output path's file extension.
 * So if the outputPath specified is "/tmp/generated_image", and the template base image is PNG, the output format will be "png".
 * If the template base image is JPG then the output path will be "/tmp/generated_image.jpg".
 *
 * Returns an error if any step fails, along with the final output path of the generated image.
 */
func GenerateImage(tmpl *Template, inputs Inputs, outputPath string) (error, string) {
	baseImg, err := LoadImageFromFile(tmpl.TemplateImage)
	if err != nil {
		return fmt.Errorf("loading base image: %v", err), ""
	}

	var canvas *image.RGBA
	if tmpl.Output.Width > 0 && tmpl.Output.Height > 0 {
		canvas = image.NewRGBA(image.Rect(0, 0, tmpl.Output.Width, tmpl.Output.Height))
		// draw scaled base to fill canvas
		scaledBase := ResizeImage(baseImg, tmpl.Output.Width, tmpl.Output.Height, ResizeModeFill)
		draw.Draw(canvas, canvas.Bounds(), scaledBase, image.Point{0, 0}, draw.Src)
		log.Printf("ImageDriver: drew scaled base image to fill canvas")
	} else {
		b := baseImg.Bounds()
		canvas = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(canvas, canvas.Bounds(), baseImg, b.Min, draw.Src)
		log.Printf("ImageDriver: drew base image")
	}

	dc := gg.NewContextForRGBA(canvas)

	// Process slots
	for _, slot := range tmpl.Slots {
		val, ok := inputs[slot.ID]
		if !ok {
			continue
		}

		if slot.IsText {
			// draw text in slot
			log.Printf("ImageDriver: drawing text for slot %s", slot.ID)
			slot.DrawTextInto(dc, val)
			continue
		}

		// load image
		imgPath := val
		img, err := LoadImageFromFile(imgPath)
		if err != nil {
			log.Printf("warning: failed to load image for slot %s: %v", slot.ID, err)
			continue
		}

		mode := slot.Mode
		if mode == "" {
			mode = ResizeModeFit
		}

		log.Printf("ImageDriver: processing image for slot %s with mode %s", slot.ID, mode)

		resized := ResizeImage(img, slot.Width, slot.Height, mode)
		// apply opacity
		finalImg := ApplyOpacity(resized, slot.Opacity)

		// If mask requested, create mask and use draw.DrawMask
		mask := MakeMask(slot.Mask, finalImg.Bounds().Dx(), finalImg.Bounds().Dy(), slot.Radius)

		// compute anchor placement
		ax := slot.AnchorX
		ay := slot.AnchorY
		if ax < 0 || ax > 1 {
			ax = 0
		}
		if ay < 0 || ay > 1 {
			ay = 0
		}
		ox := slot.X - int(float64(finalImg.Bounds().Dx())*ax)
		oy := slot.Y - int(float64(finalImg.Bounds().Dy())*ay)
		dstRect := image.Rect(ox, oy, ox+finalImg.Bounds().Dx(), oy+finalImg.Bounds().Dy())

		// prepare RGBA overlay
		rgbaOverlay := image.NewRGBA(finalImg.Bounds())
		draw.Draw(rgbaOverlay, rgbaOverlay.Bounds(), finalImg, finalImg.Bounds().Min, draw.Src)

		// draw with mask
		draw.DrawMask(canvas, dstRect, rgbaOverlay, image.Point{0, 0}, mask, image.Point{0, 0}, draw.Over)
	}

	// Save output
	outFormat := tmpl.Output.Format
	if outFormat == "" {
		ext := strings.ToLower(filepath.Ext(outputPath))
		if strings.HasPrefix(ext, ".") {
			outFormat = ext[1:]
		} else {
			outFormat = "png"
		}
	}
	outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + "." + outFormat
	log.Printf("ImageDriver: saving output image to %s with format %s", outputPath, outFormat)

	ret := SaveImageToFile(canvas, outputPath, outFormat)
	if ret != nil {
		return fmt.Errorf("saving output image: %v", ret), ""
	}

	log.Printf("Generated image saved to: %s", outputPath)

	return nil, outputPath
}
