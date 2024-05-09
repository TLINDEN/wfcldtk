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

type Superposition []*Tile

func NewTile(img image.Image, checkpoints int) (*Tile, error) {
	tile := &Tile{Image: img}
	tile.Constraints = make([]string, 4)

	id, err := GetImageHash(img)
	if err != nil {
		return nil, err
	}

	tile.Id = id

	for i, direction := range Directions {
		tile.Constraints[i] = GetConstraint(img, direction, checkpoints)
	}

	return tile, nil
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

// We consider a tile to be simple if all 4 constraints are identical
// FIXME: maybe just return the number of identical constraints?
func (tile *Tile) IsSimple() bool {
	// if all 4 constraints are  the same, the temporary map will have
	// only 1 key
	return len(map[string]int{
		tile.Constraints[North]: 0,
		tile.Constraints[East]:  0,
		tile.Constraints[South]: 0,
		tile.Constraints[West]:  0,
	}) == 1
}

func (tile *Tile) CountConstraints() int {
	return len(map[string]int{
		tile.Constraints[North]: 0,
		tile.Constraints[East]:  0,
		tile.Constraints[South]: 0,
		tile.Constraints[West]:  0,
	})
}
