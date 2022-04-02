package wfc

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

func IsImageFile(file string) bool {
	return file[len(file)-4:] == ".png"
}

func ReadDir(dir string) ([]string, error) {
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

func LoadImageFolder(folder string) ([]image.Image, error) {
	files, err := ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var images []image.Image
	for _, file := range files {
		if !IsImageFile(file) {
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
