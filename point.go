package main

import "image"

type Point image.Point

func (point Point) Shift(direction Direction) Point {
	neighbor := point

	switch direction {
	case North:
		neighbor.Y -= 1
	case East:
		neighbor.X += 1
	case South:
		neighbor.Y += 1
	case West:
		neighbor.X -= 1
	}

	return neighbor
}
