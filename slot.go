package main

// one spot in the target map,
type Slot struct {
	PossibleTiles []*Tile // starts with superposition
	Position      Point
}

func (slot Slot) Collapsed() bool {
	return len(slot.PossibleTiles) <= 1
}

func (slot Slot) Count() int {
	return len(slot.PossibleTiles)
}

// return last tile, hopefully!
func (slot Slot) GetTile() *Tile {
	return slot.PossibleTiles[0]
}
