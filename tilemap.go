package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
)

// A Tilemap is a grid of slots. Each slots holds 0-N Tiles (where the
// maximum number  of tiles  is the  "superposition").  The  width and
// height are  grid positions,  not pixel coordinates.   We use  a map
// here because it's easier to access matching Slots (neighbors etc).
type Tilemap struct {
	Width, Height int
	Slots         map[Point]*Slot
	Slotlist      []*Slot // same content, but used for iterating or sorting
	Collapsing    bool
}

// Return a new empty Tilemap
func NewTilemap(width, height int) Tilemap {
	return Tilemap{
		Width:    width,
		Height:   height,
		Slots:    make(map[Point]*Slot, width*height),
		Slotlist: make([]*Slot, width*height),
	}
}

// Put all possible tiles we have (known as "superposition") into each
// slot on the target  map. The tiles in each slot  will be later then
// reduced ("collapsed") up to the point where only 1 tile is left. At
// that point it is considered to be in collapsed state.
func (tilemap *Tilemap) Populate(superposition []*Tile) {
	pos := 0

	for y := 0; y < tilemap.Height; y++ {
		for x := 0; x < tilemap.Width; x++ {
			point := Point{X: x, Y: y}
			tilemap.Slots[point] = &Slot{PossibleTiles: superposition, Position: point}
			tilemap.Slotlist[pos] = tilemap.Slots[point]
			pos++
		}
	}
}

// Print the Tilemap (only coordinate + tile count)
func (tilemap *Tilemap) Dump() {
	tilemap.DumpData(false)
}

// Print the Tilemap (coordinate + tile constraints)
func (tilemap *Tilemap) DumpAll() {
	tilemap.DumpData(true)
}

// Actual dumper implementation
func (tilemap *Tilemap) DumpData(full bool) {
	for y := 0; y < tilemap.Height; y++ {
		for x := 0; x < tilemap.Width; x++ {
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
func (tilemap *Tilemap) Collapsed() bool {
	for _, slot := range tilemap.Slots {
		if !slot.Collapsed() {
			return false
		}
	}

	return true
}

// Return true if at least 1 slot doesn't contain a tile anymore
func (tilemap *Tilemap) Broken() bool {
	for _, slot := range tilemap.Slots {
		if slot.Broken() {
			return true
		}
	}

	return false
}

// Return  true  if  the  given  slot has  a  neighbor  in  the  given
// direction. If there's a map edge on that side, returns false.
func (tilemap *Tilemap) SlotHasNeighbor(slot *Slot, direction Direction) bool {
	point := slot.Position.MoveDirection(direction)
	//fmt.Printf("        %v => %d => %v\n", slot.Position, direction, point)
	return Exists(tilemap.Slots, point)
}

// Returns neighbor slot to the given direction, if any
func (tilemap *Tilemap) GetSlotNeighbor(slot *Slot, direction Direction) *Slot {
	point := slot.Position.MoveDirection(direction)

	if !Exists(tilemap.Slots, point) {
		log.Fatalf("no slot at position %v", point)
	}

	fmt.Printf("        returning neighbor at slot %v\n", point)
	return tilemap.Slots[point]
}

func (tilemap *Tilemap) Sort() {
	sort.Slice(tilemap.Slotlist, func(left, right int) bool {
		return tilemap.Slotlist[left].Count() < tilemap.Slotlist[right].Count()
	})
}

// Try to collapse all slots, recursively
func (tilemap *Tilemap) Collapse() error {
	for !tilemap.Collapsed() {
		tilemap.Sort()
		collapsing := false

		for _, slot := range tilemap.Slotlist {
			point := slot.Position
			fmt.Printf("looking at slot at point %v\n", point)

			if slot.Collapsed() {
				// already collapsed, ignore this time
				fmt.Println("    ignore already collapsed")
				continue
			}

			if !collapsing {
				// first slot for this round, collapse  this one
				fmt.Println("    collapsing first")
				slot.Collapse()
				collapsing = true
				continue
			}

			// if tilemap.SlotVisited(point) {
			// 	fmt.Println("    slot has already been visited")
			// 	continue
			// }

			// for  current slot, look  at each direction  and exclude
			//  any  tile  which  does  not match  one  of  the  tiles
			// of the neighbor slot.
			for _, direction := range Directions {
				fmt.Printf("    looking into direction %d\n", direction)
				if !tilemap.SlotHasNeighbor(slot, direction) {
					fmt.Printf("       slot has no neighbor in direction %d\n", direction)
					continue
				}

				neighborslot := tilemap.GetSlotNeighbor(slot, direction)

				count := slot.Count()
				slot.Exclude(neighborslot, direction)
				fmt.Printf("        reduced slot from %d to %d tiles\n", count, slot.Count())
			}
		}

		fmt.Println()

		if tilemap.Broken() {
			// FIXME: either do backtracking from here OR
			return errors.New("tilemap broken")
		}
	}

	return nil
}
