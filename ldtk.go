package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/solarlune/ldtkgo"
)

type LDTKProject struct {
	Directory string
	Project   *ldtkgo.Project
}

type TileSetSubRect struct {
	X, Y, W, H int
}

func Map2Subrect(raw map[string]any) *TileSetSubRect {
	// we need to translate this map for less typing
	return &TileSetSubRect{
		W: int(raw["w"].(float64)),
		H: int(raw["h"].(float64)),
		X: int(raw["x"].(float64)),
		Y: int(raw["y"].(float64)),
	}
}

func GetPropertyRef(entity *ldtkgo.Entity, refname string) string {
	ref := entity.PropertyByIdentifier(refname)
	if ref != nil {
		if ref.Value != nil {
			refid := ref.Value.(map[string]interface{})
			ref := refid["entityIid"].(string)

			if ref != "" {
				return ref
			}
		}
	}

	return ""
}

func GetPropertyString(entity *ldtkgo.Entity, property string) string {
	ref := entity.PropertyByIdentifier(property)
	if ref != nil {
		return ref.AsString()
	}

	return ""
}

func GetPropertyBool(entity *ldtkgo.Entity, property string) bool {
	ref := entity.PropertyByIdentifier(property)
	if ref != nil {
		return ref.AsBool()
	}

	return false
}

func GetPropertyEnum(entity *ldtkgo.Entity, property string) string {
	ref := entity.PropertyByIdentifier(property)
	if ref != nil {
		return ref.AsString()
	}

	return ""
}

func GetPropertyToggleTile(entity *ldtkgo.Entity, togglename string) *TileSetSubRect {
	ref := entity.PropertyByIdentifier(togglename)
	if ref != nil {
		return Map2Subrect(ref.AsMap())
	}

	return nil
}

// load LDTK project from disk
func LDTKLoadProjectFile(file string) LDTKProject {
	fd, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open LDTK file %s: %s", file, err)
	}
	defer fd.Close()

	fileinfo, err := fd.Stat()
	if err != nil {
		log.Fatalf("failed to stat() LDTK file %s: %s", file, err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = fd.Read(buffer)
	if err != nil {
		log.Fatalf("failed to read bytes from LDTK file %s: %s", file, err)
	}

	ldtkproject, err := ldtkgo.Read(buffer)
	if err != nil {
		panic(err)
	}

	basepath := filepath.Dir(file)

	return LDTKProject{Project: ldtkproject, Directory: basepath}
}

func LDTKGetCellsize(project LDTKProject, identifier string) int {
	level := &ldtkgo.Level{}

	for _, lvl := range project.Project.Levels {
		if lvl.Identifier == identifier {
			level = lvl
			break
		}
	}

	for _, layer := range level.Layers {
		switch layer.Type {
		case ldtkgo.LayerTypeTile:
			return layer.GridSize
		}
	}

	return 0
}

// load superposition tile array from named LDTK level
func LDTKLoadLevel(project LDTKProject, identifier string, checkpoints int) []*Tile {
	superposition := []*Tile{}

	level := &ldtkgo.Level{}

	for _, lvl := range project.Project.Levels {
		if lvl.Identifier == identifier {
			level = lvl
			break
		}
	}

	for _, layer := range level.Layers {
		switch layer.Type {
		case ldtkgo.LayerTypeTile:
			// load tile from LDTK tile layer
			tileset := layer.Tileset
			tilemap, err := Loadimage(project.Directory + "/" + tileset.Path)
			if err != nil {
				log.Fatalf("failed to load tileset %s: %s", tileset.Path, err)
			}

			if tiles := layer.AllTiles(); len(tiles) > 0 {
				for _, tileData := range tiles {
					// fetch current tile from current level from current tileset
					// FIXME: measurements are wrong!
					panic("check measurements!")
					tileimage, err := GetTileFromSpriteSheet(
						tilemap,
						tileData.Src[0],
						tileData.Src[1],
						tileData.Src[0]+layer.GridSize,
						tileData.Src[1]+layer.GridSize)
					if err != nil {
						log.Fatalf("failed to load subimage from %s: %s", tileset.Path, err)
					}

					tile := NewTile(tileimage, checkpoints)
					superposition = append(superposition, tile)

					if DEBUG {
						file := fmt.Sprintf("images/tile-debug-%d-%d.png",
							tileData.Src[0],
							tileData.Src[1])
						SavePNG(file, tileimage)
					}
				}
			}
		}
	}

	return superposition
}
