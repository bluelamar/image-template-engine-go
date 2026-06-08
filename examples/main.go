package main

import (
	"fmt"
	"log"

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
	templatePath := "test/example_template.json"
	// The inputs JSON file should provide the image and/or text content for each slot defined in the template.
	inputsPath := "test/example_input.json"
	// The output path where the generated image will be saved.
	outputPath := "/tmp/generated_image.png"

	err := iteng.ImageDriver(templatePath, inputsPath, outputPath)
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	fmt.Printf("Image generated and saved to %s\n", outputPath)
}
