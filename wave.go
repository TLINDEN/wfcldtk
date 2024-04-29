package main

import (
	"fmt"
	"image"

	"log"
)

type Wave struct {
	Output                               Tilemap
	Width, Height, Cellsize, Checkpoints int
	Superposition                        []*Tile // holds all possible tiles
}

// feed directly with tiles pre-fabricated by the caller
func NewWave(tileset image.Image, width, height, cellsize, checkpoints int) Wave {
	wave := Wave{
		Width:       width,
		Height:      height,
		Cellsize:    cellsize,
		Checkpoints: checkpoints,
		Output:      NewTilemap(width, height),
	}

	wave.SetupSuperposition(tileset)

	// FIXME: this is the point where we could pre-populate!
	wave.Output.Populate(wave.Superposition)

	return wave
}

func (wave *Wave) SetupSuperposition(tileset image.Image) {
	width := tileset.Bounds().Dx()
	height := tileset.Bounds().Dy()

	for x := 0; x < width; x += wave.Cellsize {
		for y := 0; y < height; y += wave.Cellsize {
			tileimage, err := GetTileFromSpriteSheet(tileset,
				x, y, wave.Cellsize, wave.Cellsize)
			if err != nil {
				log.Fatalf("failed to load tile image: %s\n", err)
			}

			file := fmt.Sprintf("images/tile-debug-%d-%d.png", x, y)
			SavePNG(file, tileimage)

			if !ImageIsTransparent(tileimage) {
				wave.Superposition = append(wave.Superposition,
					//NewTile(tileimage, x, y, wave.Checkpoints),
					NewTile(tileimage, wave.Checkpoints),
				)
			}
		}
	}
}
