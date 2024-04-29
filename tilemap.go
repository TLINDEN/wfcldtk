package main

import "fmt"

type Tilemap struct {
	Width, Height int
	Slots         map[Point]Slot
}

func NewTilemap(width, height int) Tilemap {
	return Tilemap{
		Width:  width,
		Height: height,
		Slots:  make(map[Point]Slot, width*height),
	}
}

func (tilemap Tilemap) Populate(superposition []*Tile) {
	for x := 0; x < tilemap.Width; x++ {
		for y := 0; y < tilemap.Height; y++ {
			point := Point{X: x, Y: y}
			tilemap.Slots[point] = Slot{PossibleTiles: superposition, Position: point}
		}
	}
}

func (tilemap Tilemap) Dump() {
	for x := 0; x < tilemap.Width; x++ {
		for y := 0; y < tilemap.Height; y++ {
			point := Point{X: x, Y: y}
			fmt.Printf("(%v):%d ", point, tilemap.Slots[point].Count())
		}
		fmt.Println()
	}
}
