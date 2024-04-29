package main

import (
	"fmt"
	"log"
)

// A Tilemap is a grid of slots. Each slots holds 0-N Tiles (where the
// maximum number  of tiles  is the  "superposition").  The  width and
// height are  grid positions,  not pixel coordinates.   We use  a map
// here because it's easier to access matching Slots (neighbors etc).
type Tilemap struct {
	Width, Height int
	Slots         map[Point]Slot
	Seenslots     map[Point]bool
}

// Return a new empty Tilemap
func NewTilemap(width, height int) Tilemap {
	return Tilemap{
		Width:     width,
		Height:    height,
		Slots:     make(map[Point]Slot, width*height),
		Seenslots: make(map[Point]bool{}, tilemap.Height*tilemap.Width),
	}
}

// Put all possible tiles we have (known as "superposition") into each
// slot on the target  map. The tiles in each slot  will be later then
// reduced ("collapsed") up to the point where only 1 tile is left. At
// that point it is considered to be in collapsed state.
func (tilemap Tilemap) Populate(superposition []*Tile) {
	for x := 0; x < tilemap.Width; x++ {
		for y := 0; y < tilemap.Height; y++ {
			point := Point{X: x, Y: y}
			tilemap.Slots[point] = Slot{PossibleTiles: superposition, Position: point}
		}
	}
}

// Print the Tilemap (only coordinate + tile count)
func (tilemap Tilemap) Dump() {
	tilemap.DumpData(false)
}

// Print the Tilemap (coordinate + tile constraints)
func (tilemap Tilemap) DumpAll() {
	tilemap.DumpData(true)
}

// Actual dumper implementation
func (tilemap Tilemap) DumpData(full bool) {
	for x := 0; x < tilemap.Width; x++ {
		for y := 0; y < tilemap.Height; y++ {
			point := Point{X: x, Y: y}
			fmt.Printf("(%v):%d", point, tilemap.Slots[point].Count())
			if full {
				fmt.Println()
				for _, tile := range tilemap.Slots[point].PossibleTiles {
					fmt.Println(tile.Dump())
				}
			}
		}
		fmt.Println()
	}
}

// Return true if  all slots are collapsed, that is  - each slots only
// contains 1 tile
func (tilemap Tilemap) Collapsed() bool {
	for _, slot := range tilemap.Slots {
		if !slot.Collapsed() {
			return false
		}
	}

	return true
}

// Return true if at least 1 slot doesn't contain a tile anymore
func (tilemap Tilemap) Broken() bool {
	for _, slot := range tilemap.Slots {
		if slot.Broken() {
			return false
		}
	}

	return false
}

// Return  true  if  the  given  slot has  a  neighbor  in  the  given
// direction. If there's a map edge on that side, returns false.
func (tilemap Tilemap) SlotHasNeighbor(slot Slot, direction Direction) bool {
	if slot.Position.X == 0 ||
		slot.Position.Y == 0 ||
		slot.Position.X == tilemap.Width ||
		slot.Position.Y == tilemap.Height {
		return false
	}

	return true
}

// Returns neighbor slot to the given direction, if any
func (tilemap Tilemap) GetSlotNeighbor(slot Slot, direction Direction) Slot {
	point := slot.Position

	switch direction {
	case North:
		point.Y--
	case South:
		point.Y++
	case East:
		point.X++
	case West:
		point.X--
	}

	if !Exists(tilemap.Slots, point) {
		log.Fatalf("no slot at position %v", point)
	}

	return tilemap.Slots[point]
}

func (tilemap Tilemap) SlotVisited(point Point) bool {
	return Exists(tilemap.Seenslots, point)
}

// Try to collapse all slots, recursively
func (tilemap Tilemap) Collapse() bool {
	for !tilemap.Collapsed() {
		tilemap.Seenslots = make(map[Point]bool, tilemap.Height*tilemap.Width)

		for point, slot := range tilemap.Slots {
			if slot.Collapsed() {
				continue
			}

			for _, direction := range Directions {
				if !tilemap.SlotHasNeighbor(slot, direction) || tilemap.SlotVisited(point) {
					continue
				}

				neighborslot := tilemap.GetSlotNeighbor(slot, direction)

				// get possible tiles matching to current slot in given direction

				// reduce current tiles based on possible tiles

				// if 1 left: superposition, continue

				tilemap.Seenslots[point] = true
			}
		}

		if tilemap.Broken() {
			log.Fatal("tilemap broken")
		}
	}

	return true
}
