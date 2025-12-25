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

func TestParseInputs(t *testing.T) {
	inputFile := "../test/test_input.json"
	inputs, err := ParseInputs(inputFile)
	if err != nil {
		t.Errorf("ParseInputs failed with error: %v", err)
	}
	if len(inputs) != 5 {
		t.Errorf("Expected 5 inputs, but got %d", len(inputs))
	}
	if inputs["motif"] != "./test/arrow_100x100.png" {
		t.Errorf("Expected motif to be './test/arrow_100x100.png', but got '%s'", inputs["motif"])
	}
	if inputs["title"] != "Introducing Text Rendering" {
		t.Errorf("Expected title to be 'Introducing Text Rendering', but got '%s'", inputs["title"])
	}
	if inputs["sub_title"] != "Flexible, Configurable & High-Quality" {
		t.Errorf("Expected sub_title to be 'Flexible, Configurable & High-Quality', but got '%s'", inputs["sub_title"])
	}
	if inputs["description"] != "This example shows how text rendering works" {
		t.Errorf("Expected description to be 'This example shows how text rendering works', but got '%s'", inputs["description"])
	}
	if inputs["footer"] != "Generated on 2025-02-15" {
		t.Errorf("Expected footer to be 'Generated on 2025-02-15', but got '%s'", inputs["footer"])
	}
}

func TestParseTemplate(t *testing.T) {
	templateFile := "../test/test_template.json"
	templ, err := ParseTemplate(templateFile)
	if err != nil {
		t.Errorf("ParseTemplate failed with error: %v", err)
	}
	if templ == nil {
		t.Errorf("Expected non-empty template, but got empty string")
	}

	if templ.TemplateImage != "test/sun_and_moon_100x100.png" {
		t.Errorf("Expected template image to be '%s', but got '%s'", templateFile, templ.TemplateImage)
	}

	output := templ.Output
	if expected := "png"; output.Format != expected {
		t.Errorf("Expected output to be '%s', but got '%s'", expected, output.Format)
	}

	width := output.Width
	if expected := 1024; width != expected {
		t.Errorf("Expected width to be %d, but got %d", expected, width)
	}

	height := output.Height
	if expected := 600; height != expected {
		t.Errorf("Expected height to be %d, but got %d", expected, height)
	}
}
