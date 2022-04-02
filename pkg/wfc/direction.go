package wfc

// Direction type, one of Up, Down, Left, Right.
//
// Used to specify the direction of a constraint between two modules.
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var Directions = []Direction{Down, Left, Right, Up}

// Opposite returns the opposite direction of "this" direction.
func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	return 0
}

// ToString returns the string representation of the direction.
func (d Direction) ToString() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	}
	return "Unknown"
}
