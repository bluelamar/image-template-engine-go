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
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

func ImageDriver(templatePath string, inputsPath string, outputPath string) error {
	tmpl, err := ParseTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err)
	}

	inputs, err := ParseInputs(inputsPath)
	if err != nil {
		return fmt.Errorf("parsing inputs: %v", err)
	}

	baseImg, err := LoadImageFromFile(tmpl.TemplateImage)
	if err != nil {
		return fmt.Errorf("loading base image: %v", err)
	}

	var canvas *image.RGBA
	if tmpl.Output.Width > 0 && tmpl.Output.Height > 0 {
		canvas = image.NewRGBA(image.Rect(0, 0, tmpl.Output.Width, tmpl.Output.Height))
		// draw scaled base to fill canvas
		scaledBase := ResizeImage(baseImg, tmpl.Output.Width, tmpl.Output.Height, ResizeModeFill)
		draw.Draw(canvas, canvas.Bounds(), scaledBase, image.Point{0, 0}, draw.Src)
	} else {
		b := baseImg.Bounds()
		canvas = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(canvas, canvas.Bounds(), baseImg, b.Min, draw.Src)
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
			slot.DrawTextInto(dc, val)
			continue
		}

		// load image
		imgPath := val
		img, err := LoadImageFromFile(imgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to load image for slot %s: %v", slot.ID, err)
			continue
		}

		mode := slot.Mode
		if mode == "" {
			mode = ResizeModeFit
		}

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

	ret := SaveImageToFile(canvas, outputPath, outFormat)
	if ret != nil {
		return fmt.Errorf("saving output image: %v", ret)
	}

	//log.Printf("Generated image saved to: %s", outputPath)

	return nil
}
