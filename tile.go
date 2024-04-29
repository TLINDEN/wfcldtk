package main

import (
	"fmt"
	"image"
)

type Tile struct {
	Type        string
	Image       image.Image
	Constraints []string // one per side
}

const (
	North = iota
	East
	South
	West
)

type Direction int

var Directions = []Direction{0, 1, 2, 3}

func NewTile(img image.Image, checkpoints int) *Tile {
	tile := &Tile{Image: img}
	tile.Constraints = make([]string, 4)

	for i, direction := range Directions {
		tile.Constraints[i] = GetConstraint(img, direction, checkpoints)
	}

	return tile
}

func (tile *Tile) Dump() string {
	return fmt.Sprintf("  [N:%s E:%s S:%s W:%s]",
		tile.Constraints[North],
		tile.Constraints[East],
		tile.Constraints[South],
		tile.Constraints[West],
	)
}
