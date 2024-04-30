package main

import (
	"fmt"
	"image"
)

type Tile struct {
	Id          string
	Type        string
	Image       image.Image
	Constraints []string // one per side
}

func NewTile(img image.Image, checkpoints int) *Tile {
	tile := &Tile{Image: img}
	tile.Constraints = make([]string, 4)
	tile.Id = GetImageHash(img)

	for i, direction := range Directions {
		tile.Constraints[i] = GetConstraint(img, direction, checkpoints)
	}

	return tile
}

func (tile *Tile) Dump() string {
	return fmt.Sprintf("  [N:%s E:%s S:%s W:%s <%s>]",
		tile.Constraints[North],
		tile.Constraints[East],
		tile.Constraints[South],
		tile.Constraints[West],
		tile.Id,
	)
}
