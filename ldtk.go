package main

import (
	"fmt"
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
func LDTKLoadProjectFile(file string) (*LDTKProject, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open LDTK file %s: %w", file, err)
	}
	defer fd.Close()

	fileinfo, err := fd.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat() LDTK file %s: %w", file, err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = fd.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes from LDTK file %s: %w", file, err)
	}

	ldtkproject, err := ldtkgo.Read(buffer)
	if err != nil {
		panic(err)
	}

	basepath := filepath.Dir(file)

	return &LDTKProject{Project: ldtkproject, Directory: basepath}, nil
}

func LDTKGetCellsize(project *LDTKProject, identifier string) int {
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
func LDTKLoadLevel(project *LDTKProject, identifier string, checkpoints int) (Superposition, error) {
	superposition := Superposition{}

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
				return nil, fmt.Errorf("failed to load tileset %s: %w", tileset.Path, err)
			}

			for _, tileData := range layer.AllTiles() {
				// fetch current tile from current level from current tileset
				tileimage, err := GetTileFromSpriteSheet(
					tilemap,
					tileData.Src[0],
					tileData.Src[1],
					layer.GridSize,
					layer.GridSize)
				if err != nil {
					return nil, fmt.Errorf("failed to load subimage from %s: %w", tileset.Path, err)
				}

				tile, err := NewTile(tileimage, checkpoints)
				if err != nil {
					return nil, nil
				}

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

	return superposition, nil
}
