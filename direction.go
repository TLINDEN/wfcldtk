package main

const (
	North = iota
	East
	South
	West
)

type Direction int

var Directions = []Direction{0, 1, 2, 3}

func GetAdverseDir(direction Direction) Direction {
	switch direction {
	case North:
		return South
	case South:
		return North
	case West:
		return East
	case East:
		return West
	}

	panic("invalid direction")
}

func (point *Point) MoveDirection(direction Direction) Point {
	newpoint := Point{point.X, point.Y}

	switch direction {
	case North:
		newpoint.Y--
	case South:
		newpoint.Y++
	case East:
		newpoint.X++
	case West:
		newpoint.X--
	}

	return newpoint
}
