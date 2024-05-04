package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"log"
)

var DEBUG bool

type Wave struct {
	OutputTilemap                        Tilemap
	Width, Height, Cellsize, Checkpoints int
	Superposition                        []*Tile // holds all possible tiles
}

// feed directly with tiles pre-fabricated by the caller
func NewWaveFromTileset(tileset image.Image,
	width, height, cellsize, checkpoints int) Wave {

	wave := Wave{
		Width:         width,
		Height:        height,
		Cellsize:      cellsize,
		Checkpoints:   checkpoints,
		OutputTilemap: NewTilemap(width, height),
	}

	wave.SetupSuperpositionTileset(tileset)

	// FIXME: this is the point where we could pre-populate!
	wave.OutputTilemap.Populate(wave.Superposition)

	return wave
}

func NewWaveFromProject(projectname, level string, width, height, checkpoints int) Wave {

	wave := Wave{
		Checkpoints: checkpoints,
		Width:       width,
		Height:      height,
	}

	project := LDTKLoadProjectFile(projectname)
	wave.Cellsize = LDTKGetCellsize(project, level)

	wave.OutputTilemap = NewTilemap(wave.Width, wave.Height)

	wave.SetupSuperpositionLDTK(project, level)

	// FIXME: this is the point where we could pre-populate!
	wave.OutputTilemap.Populate(wave.Superposition)

	return wave
}

// Create tiles  from the given  tileset, and  put all tiles  into the
// superposition, which is just a slice of all possible tiles
func (wave *Wave) SetupSuperpositionTileset(tileset image.Image) {
	width := tileset.Bounds().Dx()
	height := tileset.Bounds().Dy()

	for y := 0; y < height; y += wave.Cellsize {
		for x := 0; x < width; x += wave.Cellsize {
			tileimage, err := GetTileFromSpriteSheet(tileset,
				x, y, wave.Cellsize, wave.Cellsize)
			if err != nil {
				log.Fatalf("failed to load tile image: %s\n", err)
			}

			if DEBUG {
				file := fmt.Sprintf("images/tile-debug-%d-%d.png", x, y)
				SavePNG(file, tileimage)
			}

			if !ImageIsTransparent(tileimage) {
				wave.Superposition = append(wave.Superposition,
					//NewTile(tileimage, x, y, wave.Checkpoints),
					NewTile(tileimage, wave.Checkpoints),
				)
			}
		}
	}
}

// Same thing, but use an LDTK project file as the source
func (wave *Wave) SetupSuperpositionLDTK(project LDTKProject, level string) {
	wave.Superposition = LDTKLoadLevel(project, level, wave.Checkpoints)
}

// Collapse the wave
func (wave *Wave) Collapse(retries int) error {
	return wave.OutputTilemap.Collapse(retries)
}

func (wave *Wave) Export(filename string) error {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{wave.Width * wave.Cellsize, wave.Height * wave.Cellsize}

	renderto := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for point, slot := range wave.OutputTilemap.Slots {
		bounds := image.Rect(
			point.X*wave.Cellsize, point.Y*wave.Cellsize,
			(point.X+1)*wave.Cellsize, (point.Y+1)*wave.Cellsize,
		)

		if slot.Count() == 1 {
			tile := slot.GetTile().Image
			draw.Draw(renderto, bounds, tile, image.ZP, draw.Over)
		} else {
			red := color.RGBA{255, 0, 0, 255}
			draw.Draw(renderto, bounds, &image.Uniform{red}, image.ZP, draw.Src)
		}
	}

	return SavePNG(filename, renderto)
}
