package main

import (
	"fmt"
	"math/rand"
)

// one spot in the target map,
type Slot struct {
	PossibleTiles []*Tile // starts with superposition
	Position      Point
}

// Return true if slot is collapsed
func (slot *Slot) Collapsed() bool {
	return slot.Count() == 1
}

// Return true if slot is broken, that is it doesn't contain any tile
func (slot *Slot) Broken() bool {
	return slot.Count() == 0
}

// Tile count
func (slot *Slot) Count() int {
	return len(slot.PossibleTiles)
}

// return last tile, hopefully!
func (slot *Slot) GetTile() *Tile {
	return slot.PossibleTiles[0]
}

func (slot *Slot) Collapse() {
	tile := slot.PossibleTiles[rand.Intn(slot.Count())]
	slot.PossibleTiles = []*Tile{tile}
}

/*
Check  all tiles  on current  slot. If  the side  of a  tile pointing
towards one of  the neighbor slot's tiles matching sides,  it will be
kept. So we iterate over all current tiles and all neighbor tiles and
look which ones adversial sides match.

	 ----        ----
	|    |  =>  |    |
	 ----        ----

slot.dir     other.adverse

	east   =>  west
*/
func (slot *Slot) Exclude(otherslot *Slot, direction Direction) {
	adversedirection := GetAdverseDir(direction)

	keeptiles := map[string]*Tile{}

	for _, othertile := range otherslot.PossibleTiles {
		for _, tile := range slot.PossibleTiles {
			if tile.Constraints[direction] == othertile.Constraints[adversedirection] {
				if !Exists(keeptiles, tile.Id) {
					fmt.Printf("        matching this dir %d %s <=> other dir %d %s\n",
						direction, tile.Constraints[direction], adversedirection, othertile.Constraints[adversedirection])

					keeptiles[tile.Id] = tile
				}
			}
		}
	}

	newtiles := make([]*Tile, len(keeptiles))
	i := 0
	for checksum := range keeptiles {
		fmt.Printf("            %s\n", checksum)
		newtiles[i] = keeptiles[checksum]

		i++
	}

	slot.PossibleTiles = newtiles
}
