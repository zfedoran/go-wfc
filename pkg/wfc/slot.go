package wfc

import (
	"math/rand"
)

// A Slot is a single point in space inside the Wave function. It can be in a
// superposition of many modules until collapsed.
//
// A Slot that has no superposition (length == 0) is considered to be in a
// contridiction state. Meaning that the wave function was not able to resolve it.
//
// A Slot that has a single module in its superposition is considered to be a
// collapsed slot and only has one possible module at its coordinates.
type Slot struct {
	X, Y          int       // Coordinates of the slot
	Superposition []*Module // Possible modules at the slot
}

// Collapse chooses a random module from the list of superpositions available to
// it. Its superposition list is set to the single module chosen.
func (s *Slot) Collapse() {
	module := s.Superposition[rand.Intn(len(s.Superposition))]
	s.Superposition = []*Module{module}
}

// IsPossibleFunc is a function that returns whether or not a module is possible
// given a slot and direction. Use this if you'd like custom logic.
type IsPossibleFunc func(state *Module, from, to *Slot, d Direction) bool

// DefaultIsPossibleFunc returns whether or not a module is possible given a
// slot and direction.
func DefaultIsPossibleFunc(state *Module, from, to *Slot, d Direction) bool {
	return state.IsPossibleFrom(from, d)
}
