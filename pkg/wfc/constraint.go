package wfc

import (
	"crypto/sha256"
	"fmt"
	"image"
)

// Adjacency constraint type.
//
// This is a an arbitrary value that should be the same between two or more
// modules. It represents modules that can be next to each other.
//
// This codebase combines 3 color values along the edges of input tiles and then
// uses sha256 to generate a hash. Only the first 8 bytes are kept.
//
// The color values are reduced to 4 bits each, discarding the least significant
// bits. This rounds the color values to allow for some flexibility.
type ConstraintId [8]byte

// ConstraintFunc is a function that returns an adjacency hash for an image tile
// in a specified direction.
type ConstraintFunc func(image.Image, Direction) ConstraintId

// The default constraint function uses color values to generate an adjacency.
var DefaultConstraintFunc ConstraintFunc = GetConstraintId

// GetConstraintId returns the adjacency constraint id for the given tile image
// in the provided direction.
func GetConstraintId(img image.Image, d Direction) ConstraintId {
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	// Divide into 4 chunks and nudge slightly to not align on popular grids
	u := w / 4.0 - (img.Bounds().Max.X/10.0)
	v := h / 4.0 - (img.Bounds().Max.Y/10.0)

	var hash string
	var a, b, c Color

	switch d {
	case Up:
		a = GetColor(img, u*1, 0)
		b = GetColor(img, u*2, 0)
		c = GetColor(img, u*3, 0)
	case Down:
		a = GetColor(img, u*1, h-1)
		b = GetColor(img, u*2, h-1)
		c = GetColor(img, u*3, h-1)
	case Left:
		a = GetColor(img, 0, v*1)
		b = GetColor(img, 0, v*2)
		c = GetColor(img, 0, v*3)
	case Right:
		a = GetColor(img, w-1, v*1)
		b = GetColor(img, w-1, v*2)
		c = GetColor(img, w-1, v*3)
	}

	// Generate a hash from the colors
	hash = HexFromColor(DiscardLeastSignificantBits(a, 2)) + ":" +
		HexFromColor(DiscardLeastSignificantBits(b, 2)) + ":" +
		HexFromColor(DiscardLeastSignificantBits(c, 2))

	sum := sha256.Sum256([]byte(hash))
	res := fmt.Sprintf("%x", sum)[:8]

	var id ConstraintId
	copy(id[:], res)
	return id
}

// GetConstraintFromHex returns the adjacency constraint id for the given hex
func GetConstraintFromHex(s string) ConstraintId {
	var id ConstraintId
	copy(id[:], s)
	return id
}

type Color [4]uint8

func GetColor(img image.Image, x, y int) Color {
	r, g, b, a := img.At(x, y).RGBA()
	return Color{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func DiscardLeastSignificantBits(c Color, bits int) Color {
	return Color{c[0] >> bits, c[1] >> bits, c[2] >> bits, c[3] >> bits}
}

func HexFromColor(c Color) string {
	return fmt.Sprintf("%02x%02x%02x%02x", c[0], c[1], c[2], c[3])
}
