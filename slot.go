package main

import (
	"fmt"
	"math/rand"
)

// one spot in the target map,
type Slot struct {
	PossibleTiles         Superposition // starts with superposition
	PreviousPossibleTiles Superposition // backup
	Position              Point
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
	slot.PossibleTiles = Superposition{tile}
}

func (slot *Slot) Copy() {
	slot.PreviousPossibleTiles = slot.PossibleTiles
}

func (slot *Slot) Backtrack() {
	slot.PossibleTiles = slot.PreviousPossibleTiles
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
func (slot *Slot) Exclude(otherslot *Slot, direction Direction) map[string]*Tile {
	adversedirection := GetAdverseDir(direction)

	keeptiles := map[string]*Tile{}

	for _, othertile := range otherslot.PossibleTiles {
		for _, tile := range slot.PossibleTiles {
			if tile.Constraints[direction] == othertile.Constraints[adversedirection] {
				if !Exists(keeptiles, tile.Id) {
					if DEBUG {
						fmt.Printf("        matching this dir %d %s <=> other dir %d %s\n",
							direction, tile.Constraints[direction], adversedirection,
							othertile.Constraints[adversedirection])
					}

					keeptiles[tile.Id] = tile
				}
			}
		}
	}

	return keeptiles
}

// Collapse    possible   tiles    on    this    slot   by    neighbor
// constraints.  Neighbors  are  given  in the  slice  arg,  index  ==
// direction, empty slice item means no neighbor in that direction
func (slot *Slot) CollapseByConstraints(neighbors []*Slot) {
	neighborcount := 0
	tileregistry := map[string]*Tile{}
	tilecounter := map[string]int{}

	// check all neighbor slots
	for direction, otherslot := range neighbors {
		if otherslot == nil {
			continue
		}

		// register how many neighbors there are
		neighborcount++

		//  for  every registered  neighbor in  the given  direction
		// fetch all tiles matching the neighbor constraint
		tiles := slot.Exclude(otherslot, Direction(direction))

		// count each matching tile
		for id, tile := range tiles {
			tileregistry[id] = tile
			tilecounter[id]++
		}
	}

	newtiles := Superposition{}

	// check  which tiles  are possible matches  on ALL  neighbors and
	// register them
	for id, count := range tilecounter {
		if count == neighborcount {
			newtiles = append(newtiles, tileregistry[id])
		}
	}

	// set possible tiles to the new reduced tile slice
	slot.PossibleTiles = newtiles
}
