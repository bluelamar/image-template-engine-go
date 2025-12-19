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
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	imagedraw "golang.org/x/image/draw"
)

// LoadImageFromFile decodes common image formats
func LoadImageFromFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// SaveImageToFile encodes and saves image to file with specified format
func SaveImageToFile(img image.Image, outpath, format string) error {
	out, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer out.Close()

	format = strings.ToLower(format)
	switch format {
	case "png":
		enc := png.Encoder{CompressionLevel: png.BestCompression}
		return enc.Encode(out, img)
	case "jpg", "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 92})
	case "gif":
		return gif.Encode(out, img, nil)
	default:
		// fallback to png
		enc := png.Encoder{CompressionLevel: png.BestCompression}
		return enc.Encode(out, img)
	}
}

// ResizeImage implements fill/fit/cover.
func ResizeImage(src image.Image, dstW, dstH int, mode ResizeMode) image.Image {
	if dstW <= 0 || dstH <= 0 {
		return src
	}
	srcW := src.Bounds().Dx()
	srcH := src.Bounds().Dy()

	if srcW == 0 || srcH == 0 {
		return src
	}

	switch mode {
	case ResizeModeFill:
		// direct stretch
		dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
		imagedraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), imagedraw.Over, nil)
		return dst
	case ResizeModeFit:
		scale := minf(float64(dstW)/float64(srcW), float64(dstH)/float64(srcH))
		nw := int(float64(srcW) * scale)
		nh := int(float64(srcH) * scale)
		dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
		imagedraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), imagedraw.Over, nil)
		return dst
	case ResizeModeCover:
		scale := maxf(float64(dstW)/float64(srcW), float64(dstH)/float64(srcH))
		nw := int(float64(srcW) * scale)
		nh := int(float64(srcH) * scale)
		d := image.NewRGBA(image.Rect(0, 0, nw, nh))
		imagedraw.CatmullRom.Scale(d, d.Bounds(), src, src.Bounds(), imagedraw.Over, nil)
		// center-crop to dstW x dstH
		x0 := (nw - dstW) / 2
		y0 := (nh - dstH) / 2
		out := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
		draw.Draw(out, out.Bounds(), d, image.Point{x0, y0}, draw.Src)
		return out
	default:
		// default to fit
		scale := minf(float64(dstW)/float64(srcW), float64(dstH)/float64(srcH))
		nw := int(float64(srcW) * scale)
		nh := int(float64(srcH) * scale)
		dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
		imagedraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), imagedraw.Over, nil)
		return dst
	}
}

func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// ApplyOpacity multiplies alpha channel by opacity
func ApplyOpacity(img image.Image, opacity float64) image.Image {
	// log.Printf("applyOpacity: opacity: %f", opacity)
	if opacity >= 0.9999 {
		return img
	}
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bb, a := rgba.At(x, y).RGBA()
			alpha := uint8((float64(a>>8) * opacity))
			rgba.SetRGBA(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(bb >> 8), alpha})
		}
	}
	return rgba
}

// MakeMask returns an *image.Alpha mask for the requested shape
func MakeMask(maskType string, w, h int, radius float64) *image.Alpha {
	if maskType == "" {
		m := image.NewAlpha(image.Rect(0, 0, w, h))
		// fill opaque
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				m.SetAlpha(x, y, color.Alpha{255})
			}
		}
		return m
	}

	dc := gg.NewContext(w, h)
	dc.Clear()
	dc.SetRGBA(0, 0, 0, 1)

	if maskType == "circle" {
		cx := float64(w) / 2
		cy := float64(h) / 2
		r := float64(min(w, h)) / 2
		dc.DrawCircle(cx, cy, r)
		dc.Fill()
	} else if maskType == "rounded" {
		r := radius
		if r <= 0 {
			r = float64(min(w, h)) * 0.12
		}
		dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), r)
		dc.Fill()
	} else {
		dc.DrawRectangle(0, 0, float64(w), float64(h))
		dc.Fill()
	}

	rgba := dc.Image().(*image.RGBA)
	alpha := image.NewAlpha(rgba.Bounds())
	for y := 0; y < rgba.Bounds().Dy(); y++ {
		for x := 0; x < rgba.Bounds().Dx(); x++ {
			_, _, _, a := rgba.At(x, y).RGBA()
			alpha.SetAlpha(x, y, color.Alpha{uint8(a >> 8)})
		}
	}
	return alpha
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DrawTextInto draws text into the canvas using gg and supports wrapping and alignment
// Note: dc must be initialized with the correct size before calling this function
// TODO: support vertical alignment
// Note: gg.AlignLeft is used for simplicity; adjust as needed for other alignments
// TODO: support more text options like line spacing, etc.
// TODO: support loading fonts other than from filesystem
func (slot Slot) DrawTextInto(dc *gg.Context, text string) {
	// log.Printf("drawTextInto: slot: %v: text: %s", slot, text)
	// Load font if provided
	var err error
	if slot.TextOpts.FontPath != "" {
		err = dc.LoadFontFace(slot.TextOpts.FontPath, slot.TextOpts.FontSize)
	} else {
		// attempt to load default system font; if fails, gg has a builtin fallback
		// FIX err = dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", slot.TextOpts.FontSize)
		ttfFile := os.Getenv("FONT_TTF")
		ttfPath := filepath.Join(os.Getenv("FONT_DIR"), ttfFile)
		err = dc.LoadFontFace(ttfPath, slot.TextOpts.FontSize)
	}

	if err != nil {
		log.Printf("warning: failed to load font face: %v", err)
	}

	// parse color
	if slot.TextOpts.Color != "" {
		// assume hex: #RRGGBB : slot.TextOpts.Color
		dc.SetHexColor(slot.TextOpts.Color)
	} else {
		dc.SetRGB(0, 0, 0)
	}

	// compute anchor point inside slot
	ax := slot.AnchorX
	ay := slot.AnchorY
	if ax < 0 || ax > 1 {
		ax = 0.0
	}
	if ay < 0 || ay > 1 {
		ay = 0.0
	}
	px := float64(slot.X) + float64(slot.Width)*ax
	py := float64(slot.Y) + float64(slot.Height)*ay

	// horizontal alignment mapping
	hAlign := strings.ToLower(slot.TextOpts.AlignX)
	var anchorX float64
	switch hAlign {
	case "center", "centre":
		anchorX = 0.5
	case "right":
		anchorX = 1.0
	default:
		anchorX = 0.0
	}

	// vertical alignment mapping
	vAlign := strings.ToLower(slot.TextOpts.AlignY)
	var anchorY float64
	switch vAlign {
	case "middle", "center":
		anchorY = 0.5
	case "bottom":
		anchorY = 1.0
	default:
		anchorY = 0.0
	}

	// wrapped or single-line
	if slot.TextOpts.Wrap && slot.TextOpts.MaxWidth > 0 {
		// DrawStringWrapped(x, y, ax, ay, width, lineSpacing, align)
		dc.DrawStringWrapped(text, px, py, anchorX, anchorY, float64(slot.TextOpts.MaxWidth), 1.4, gg.AlignLeft)
	} else {
		dc.DrawStringAnchored(text, px, py, anchorX, anchorY)
	}
}
