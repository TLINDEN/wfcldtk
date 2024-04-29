package main

import "math/rand"

// one spot in the target map,
type Slot struct {
	PossibleTiles []*Tile // starts with superposition
	Position      Point
}

// Return true if slot is collapsed
func (slot Slot) Collapsed() bool {
	return slot.Count() == 1
}

// Return true if slot is broken, that is it doesn't contain any tile
func (slot Slot) Broken() bool {
	return slot.Count() == 0
}

// Tile count
func (slot Slot) Count() int {
	return len(slot.PossibleTiles)
}

// return last tile, hopefully!
func (slot Slot) GetTile() *Tile {
	return slot.PossibleTiles[0]
}

func (slot Slot) Collapse() {
	tile := slot.PossibleTiles[rand.Intn(slot.Count())]
	slot.PossibleTiles = []*Tile{tile}
}
