package wfc

import (
	"image"
)

// Module represents a single module in the wave function as described by Oskar
// Stalberg. A module is a possible tile that might exist at a slot in the wave
// function grid. It can be thought of as a single state of a superposition.
type Module struct {
	Index       int             // The index of the module in the input tiles
	Adjacencies [4]ConstraintId // Adjacency constraints for each direction
	Image       image.Image     // The tile image for the module
}

// IsPossibleFrom returns true if the given module is possible from the given
// direction.
func (m *Module) IsPossibleFrom(from *Slot, forward Direction) bool {
	backward := forward.Opposite()

	for _, c := range from.Superposition {
		if m.Adjacencies[backward].Equal(c.Adjacencies[forward]) {
			return true
		}
	}

	return false
}
