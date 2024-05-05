package main

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type Stats struct {
	Superpositions, Backtracked, Rounds int
	Duration                            time.Duration
	RoundsDuration                      []time.Duration
}

// A Tilemap is a grid of slots. Each slots holds 0-N Tiles (where the
// maximum number  of tiles  is the  "superposition").  The  width and
// height are  grid positions,  not pixel coordinates.   We use  a map
// here because it's easier to access matching Slots (neighbors etc).
type Tilemap struct {
	Width, Height int
	Slots         map[Point]*Slot
	Slotlist      []*Slot // same content, but used for iterating or sorting
	Copylist      []*Slot
	Collapsing    bool
	Stats         Stats
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
func (tilemap *Tilemap) Populate(superposition Superposition) {
	pos := 0

	for y := 0; y < tilemap.Height; y++ {
		for x := 0; x < tilemap.Width; x++ {
			point := Point{X: x, Y: y}
			tilemap.Slots[point] = &Slot{PossibleTiles: superposition, Position: point}
			tilemap.Slotlist[pos] = tilemap.Slots[point]
			pos++
		}
	}

	tilemap.Stats.Superpositions = len(superposition)
}

// Make a copy of the current possibility space for backtracking
func (tilemap *Tilemap) Copy() {
	for _, slot := range tilemap.Slots {
		slot.Copy()
	}
}

// Restore previous tile set, thus backtrack one step
func (tilemap *Tilemap) Backtrack() {
	for _, slot := range tilemap.Slots {
		slot.Backtrack()
	}

	tilemap.Stats.Backtracked++
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
func (tilemap *Tilemap) GetSlotNeighbor(slot *Slot, direction Direction) (*Slot, error) {
	point := slot.Position.MoveDirection(direction)

	if !Exists(tilemap.Slots, point) {
		return nil, fmt.Errorf("no slot at position %v", point)
	}

	if DEBUG {
		fmt.Printf("        returning neighbor at slot %v\n", point)
	}
	return tilemap.Slots[point], nil
}

// Sort helper, sort Slot slice by tile count, lowest count goes first
func (tilemap *Tilemap) Sort() {
	sort.Slice(tilemap.Slotlist, func(left, right int) bool {
		return tilemap.Slotlist[left].Count() < tilemap.Slotlist[right].Count()
	})
}

// Try to collapse all slots, recursively
func (tilemap *Tilemap) Collapse(retries int) error {
	tries := 0

	for !tilemap.Collapsed() {
		start := time.Now()

		//  make a  backup of  the current  state of  the tilemap.  If
		//  collapsing   fails,  we  can   restore  it  and   thus  do
		// backtracking.
		tilemap.Copy()

		// we sort the slots by tile count, that way the slot with the
		// lowest entropy  goes first, which we collapse  at the start
		// of every loop run.
		tilemap.Sort()

		// only collapse 1 slot per run
		collapsing := false

		for _, slot := range tilemap.Slotlist {
			point := slot.Position
			if DEBUG {
				fmt.Printf("looking at slot at point %v\n", point)
			}

			if slot.Collapsed() {
				// already collapsed, ignore this time
				if DEBUG {
					fmt.Println("    ignore already collapsed")
				}
				continue
			}

			if !collapsing {
				// first slot for this round, collapse  this one
				if DEBUG {
					fmt.Println("    collapsing first")
				}
				slot.Collapse()
				collapsing = true
				continue
			}

			neighbors := make([]*Slot, 4)

			// for  current slot, look  at each direction  and exclude
			//  any  tile  which  does  not match  one  of  the  tiles
			// of the neighbor slot.
			for _, direction := range Directions {
				if DEBUG {
					fmt.Printf("    looking into direction %d\n", direction)
				}
				if !tilemap.SlotHasNeighbor(slot, direction) {
					if DEBUG {
						fmt.Printf("       slot has no neighbor in direction %d\n", direction)
					}
					continue
				}

				neighborslot, err := tilemap.GetSlotNeighbor(slot, direction)
				if err != nil {
					return err
				}

				neighbors[direction] = neighborslot

			}
			slot.CollapseByConstraints(neighbors)
		}

		if DEBUG {
			fmt.Println()
		}

		if tilemap.Broken() {
			if tries < retries {
				if DEBUG {
					fmt.Println("BACKTRACKING")
				}
				tilemap.Backtrack()

				//tilemap.Printstats()
				fmt.Printf("tries: %d, retries: %d\n", tries, retries)
				tries++
			} else {
				fmt.Printf("error tries: %d, retries: %d\n", tries, retries)
				return errors.New("tilemap broken too many times")
			}
		}

		elapsed := time.Since(start)
		tilemap.Stats.Rounds++
		tilemap.Stats.RoundsDuration = append(tilemap.Stats.RoundsDuration, elapsed)

	}

	fmt.Printf("Collapsed: %t, Broken: %t\n", tilemap.Collapsed(), tilemap.Broken())

	return nil
}

func (tilemap *Tilemap) Printstats() {
	for _, dur := range tilemap.Stats.RoundsDuration {
		tilemap.Stats.Duration += dur
	}

	fmt.Printf("Superpositions: %d\n", tilemap.Stats.Superpositions)
	fmt.Printf("         Slots: %d\n", len(tilemap.Slots))
	fmt.Printf("        Rounds: %d\n", tilemap.Stats.Rounds)
	fmt.Printf("   Backtracked: %d\n", tilemap.Stats.Backtracked)
	fmt.Printf("    time taken: %s\n", tilemap.Stats.Duration)
}
