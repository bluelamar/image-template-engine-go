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
	"encoding/json"
	"os"
)

// Inputs map slotID -> image path or text
// The Inputs should match the needed slots in the Template
type Inputs map[string]string

// ResizeMode options
// The option can be specified in the Slot struct as Mode field
type ResizeMode string

// fill - stretch the image to fit the slot
// fit - shrink the image to fit the slot
// cover - shrink the image to cover the slot
const (
	ResizeModeFill  ResizeMode = "fill"
	ResizeModeFit   ResizeMode = "fit"
	ResizeModeCover ResizeMode = "cover"
)

// Slot defines either an image or text placement in the base image
type Slot struct {
	ID       string     `json:"id"`
	X        int        `json:"x"`
	Y        int        `json:"y"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Mask     string     `json:"mask,omitempty"`     // circle, rounded, or empty
	Radius   float64    `json:"radius,omitempty"`   // for rounded
	AnchorX  float64    `json:"anchor_x,omitempty"` // 0..1
	AnchorY  float64    `json:"anchor_y,omitempty"` // 0..1
	Mode     ResizeMode `json:"mode,omitempty"`     // ResizeMode: fill/fit/cover
	Opacity  float64    `json:"opacity,omitempty"`  // 0.0 - 1.0
	IsText   bool       `json:"is_text,omitempty"`
	TextOpts TextOpt    `json:"text_opts,omitempty"`
}

// TextOpt defines text options for a Slot
type TextOpt struct {
	FontPath   string  `json:"font_path,omitempty"`   // filesystem path
	FontName   string  `json:"font_name,omitempty"`   // system font name (e.g., "Arial", "Helvetica")
	FontSource string  `json:"font_source,omitempty"` // "file", "system", "url", "embedded", or "" for auto
	FontURL    string  `json:"font_url,omitempty"`    // URL to download font from
	FontSize   float64 `json:"font_size,omitempty"`
	Color      string  `json:"color,omitempty"`   // hex like #RRGGBB
	AlignX     string  `json:"align_x,omitempty"` // left, center, right
	AlignY     string  `json:"align_y,omitempty"` // top, middle, bottom
	Wrap       bool    `json:"wrap,omitempty"`
	MaxWidth   int     `json:"max_width,omitempty"` // px for wrapping
}

// Template defines the base image, Output options, and Slots
type Template struct {
	// TemplateImage is the path to the base image to use for the template
	TemplateImage string `json:"template_image"`
	Output        Output `json:"output"`
	Slots         []Slot `json:"slots"`
}

// Output defines the output image size and format
type Output struct {
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Format string `json:"format,omitempty"` // png, jpg, gif
}

// ParseTemplate reads and parses the JSON Template file
func ParseTemplate(path string) (*Template, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t Template
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// ParseInputs reads and parses the JSON Inputs file
func ParseInputs(path string) (Inputs, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var in Inputs
	if err := json.Unmarshal(b, &in); err != nil {
		return nil, err
	}
	return in, nil
}
