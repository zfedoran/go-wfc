package wfc

import (
	"bytes"
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

// Equal returns true if the two adjacency constraints are equal.
func (c ConstraintId) Equal(o ConstraintId) bool {
	return bytes.Equal(c[:], o[:])
}

// ConstraintFunc is a function that returns an adjacency hash for an image tile
// in a specified direction.
type ConstraintFunc func(image.Image, Direction) ConstraintId

// The default constraint function uses color values to generate an adjacency.
var DefaultConstraintFunc ConstraintFunc = GetConstraintFunc(3)

// GetConstraintFunc returns a constraint function that uses the given number of
// color lookups
func GetConstraintFunc(count int) ConstraintFunc {
	count += 1
	return func(img image.Image, dr Direction) ConstraintId {
		// returns the adjacency constraint id for the given tile image
		// in the provided direction.

		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y

		u := w / count
		v := h / count

		var hash string
		points := make([]Color, count)

		for i := 0; i < count; i++ {
			switch dr {
			case Up:
				points[i] = GetColor(img, u+i*u, 0)
			case Down:
				points[i] = GetColor(img, u+i*u, h-1)
			case Left:
				points[i] = GetColor(img, 0, v+i*v)
			case Right:
				points[i] = GetColor(img, w-1, v+i*v)
			}
		}

		// Generate a hash from the colors
		hash = ""
		for _, c := range points {
			hash += HexFromColor(c)
		}

		sum := sha256.Sum256([]byte(hash))
		res := fmt.Sprintf("%x", sum)[:8]

		var id ConstraintId
		copy(id[:], res)
		return id
	}
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
