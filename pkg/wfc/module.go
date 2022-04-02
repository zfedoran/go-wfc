package wfc

import (
	"image"
)

type Module struct {
	Adjacencies [4]ConstraintId
	Image       image.Image
}

func (m *Module) IsPossibleFrom(from *Slot, forward Direction) bool {
	backward := forward.Opposite()

	for _, c := range from.Superposition {
		if m.Adjacencies[backward] == c.Adjacencies[forward] {
			return true
		}
	}

	return false
}
