package main

import (
	"fmt"
	"time"

	"github.com/zfedoran/go-wfc/pkg/wfc"
)

// Use 2 color lookups to generate adjacency hash values
var constraintFn = wfc.GetConstraintFunc(3)

func collapseWave(tileset_folder, output_image string) {
	// This is just a `[]image.Image`, you can use whatever loader function you'd like
	images, err := wfc.LoadImageFolder(tileset_folder)
	if err != nil {
		panic(err)
	}

	// The random seed to use when collapsing the wave
	// (given the same seed number, the Collapse() fn would generate the same state every time)
	seed := int(time.Now().UnixNano())

	width := 40
	height := 5

	// Setup the initialized state
	wave := wfc.NewWithCustomConstraints(images, width, height, constraintFn)
	wave.Initialize(seed)

	// Force the top tiles to be sky/void
	sky := wfc.GetConstraintFromHex("d4789c1e")
	for i := 0; i < width; i++ {
		slot := wave.PossibilitySpace[i]
		modules := make([]*wfc.Module, 0)
		for _, m := range slot.Superposition {
			if m.Adjacencies[wfc.Up] == sky {
				modules = append(modules, m)
			}
		}
		slot.Superposition = modules
	}

	// Force bottom tiles to be ground
	ground := wfc.GetConstraintFromHex("7ed0cfd4")
	for i := (height - 1) * width; i < width + (height - 1) * width; i++ {
		slot := wave.PossibilitySpace[i]
		modules := make([]*wfc.Module, 0)
		for _, m := range slot.Superposition {
			if m.Adjacencies[wfc.Left] == ground || m.Adjacencies[wfc.Right] == ground {
				modules = append(modules, m)
			}
		}
		slot.Superposition = modules
	}

	// Collapse the wave function (make up to 100 attempts)
	err = wave.Collapse(200)
	if err != nil {
		// don't panic here, we want to generate the image anyway
		fmt.Printf("unable to generate: %v", err)
	}

	// Lets generate an image
	output := wave.ExportImage()
	output_file := fmt.Sprintf("%d%s", seed, output_image)

	wfc.SaveImage(output_file, output)
	fmt.Printf("Image saved to: %s\n", output_file)
}

func printAdjacencyHashValues(input_tileset string) {
	fmt.Printf("Adjacency hash values:\n\n")

	images, err := wfc.LoadImageFolder(input_tileset)
	if err != nil {
		panic(err)
	}

	// We could use pretty table to do this, but this is just a demo and I don't
	// want the extra dependency.

	fmt.Println("|-------|----------|----------|")
	fmt.Println("|Tile\t|Direction |Hash      |")
	fmt.Println("|-------|----------|----------|")
	for i, img := range images {
		for _, d := range wfc.Directions {
			fmt.Printf("|%d\t|%s\t   | %s | %dx%d\n", i, d.ToString(), constraintFn(img, d), img.Bounds().Max.X, img.Bounds().Max.Y)
		}
		fmt.Printf("|- - - -|- - - - - |- - - - - |\n")
	}
	fmt.Printf("|-------|----------|----------|\n\n")
}

func main() {
	input_tileset := "./assets/"
	output_image_suffix := "_out.png"

	// Print the adjacency hash values for the provided tileset.
	printAdjacencyHashValues(input_tileset)

	// Generate an image from the tileset.
	collapseWave(input_tileset, output_image_suffix)


}