package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bluelamar/image-template-engine-go/iteng"
)

/*
This example demonstrates how to use the iteng package to generate an image
based on a template and input data. It loads a template JSON file that defines
the layout and slots, an inputs JSON file that provides the content for those
slots, and generates the final image saved to the specified output path.
*/
func main() {
	// The template JSON file should define the base image, output options, and slots with text options.
	templatePath := os.Getenv("ITENG_EX_TEMPL")
	if templatePath == "" {
		templatePath = "testdata/example_template.json"
	}

	log.Printf("Template file=%s\n", templatePath)

	// The inputs JSON file should provide the image and/or text content for each slot defined in the template.
	inputsPath := os.Getenv("ITENG_EX_INPUTS")
	if inputsPath == "" {
		inputsPath = "testdata/example_input.json"
	}

	log.Printf("Inputs file=%s\n", inputsPath)

	// The output path where the generated image will be saved.
	// Note that the output format (e.g., PNG, JPEG) is determined by the template's output settings.
	// So the final output file path returned would be "/tmp/generated_image.png" or "/tmp/generated_image.jpg" depending on the template's output format.
	outputPath := "/tmp/generated_image"

	log.Printf("Output path=%s\n", outputPath)

	err, outputPath := iteng.ImageDriver(templatePath, inputsPath, outputPath)
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	fmt.Printf("Image generated and saved to %s\n", outputPath)
}
