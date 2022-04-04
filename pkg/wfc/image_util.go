package wfc

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)

// LoadImageFolder loads all images in a directory and returns them as a slice.
func LoadImageFolder(folder string) ([]image.Image, error) {
	files, err := readDir(folder)
	if err != nil {
		return nil, err
	}

	var images []image.Image
	for _, file := range files {
		if !isImageFile(file) {
			continue
		}
		img, err := LoadImage(folder + "/" + file)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}

// LoadImage loads an image from a file path.
func LoadImage(file string) (image.Image, error) {
	raw, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer raw.Close()
	img, _, err := image.Decode(raw)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// SaveImage saves an image to a file.
func SaveImage(file string, img image.Image) error {
	// Check if file exists, remove if it does
	if _, err := os.Stat(file); err == nil {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	// Save new file
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	png.Encode(f, img)
	return nil
}

// GetTileFromSpriteSheet returns a tile from a sprite sheet.
func GetTileFromSpriteSheet(img image.Image, x, y, width, height int) (image.Image, error) {
	// Create new image
	outputImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Copy pixels
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			outputImg.Set(i, j, img.At(x*width+i, y*height+j))
		}
	}

	return outputImg, nil
}

// readDir reads a directory and returns a slice of file names.
func readDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, file := range files {
		names = append(names, file.Name())
	}

	return names, nil
}

// isImageFile checks if a file is an image file.
func isImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpeg", ".jpg", ".png", ".gif":
		return true
	default:
		return false
	}
}
