package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Tileset struct {
	Identifier   string
	TileGridSize float64
	RelPath      string
}

// Flip bits  - first bit is for X-flip,
// second is  for Y. 0 = no  flip, 1 = horizontal flip,  2 = vertical
// flip, 3 = both flipped

type Tile struct {
	Position []int `json:"px"`  // pixel position on target tileset
	Src      []int `json:"src"` // pixel position of tile from source tileset
	Flip     byte  `json:"f"`
	ID       int   `json:"t"`
}

type Layer struct {
	Iid            string
	LayerDefUid    int     `json:"layerDefUid"`
	LevelId        int     `json:"levelId"`
	Seed           int     `json:"seed"`
	Visible        bool    `json:"visible"`
	Width          int     `json:"__cWid"` // layer width in points
	Height         int     `json:"__cHei"` // layer height in points
	GridSize       int     `json:"__gridSize"`
	Identifier     string  `json:"__identifier"`
	Opacity        int     `json:"__opacity"`
	TilesetDefUid  int     `json:"__tilesetDefUid"`
	TilesetRelPath string  `json:"__tilesetRelPath"`
	Tiles          []*Tile `json:"gridTiles"`
}

type Level struct {
	BGColor         string  `json:"bgColor"`
	BGPivotX        float64 `json:"bgPivotX"`
	BGPivotY        float64 `json:"bgPivotY"`
	BGPos           float64 `json:"bgPos"`
	BGRelPath       string  `json:"bgRelPath"`
	ExternalRelPath string  `json:"externalRelPath"`
	Identifier      string
	Iid             string
	Layers          []Layer `json:"layerInstances"`
}

func main() {
	data := readfile("demo.ldtk")

	// read the whole JSON
	jsonstr := string(data)

	// extract just the tileset
	value := gjson.Get(jsonstr, "defs.tilesets")
	var tilesets []Tileset
	err := json.Unmarshal([]byte(value.Raw), &tilesets)
	if err != nil {
		panic(err)
	}

	//repr.Println(tilesets)

	// extract the level data
	levelval := gjson.Get(jsonstr, "levels")
	var levels []Level
	err = json.Unmarshal([]byte(levelval.Raw), &levels)
	if err != nil {
		panic(err)
	}

	// create a new level
	mylevel := Level{
		BGPivotX:   0.5,
		BGPivotY:   0.5,
		Identifier: "testing",
		Iid:        "4502d4bc-38dd-4903-85c6-bcef7ad90208",
		Layers: []Layer{
			{
				Iid:            "78a5edb8-176c-4af6-b130-38deac0047db",
				LayerDefUid:    2,
				LevelId:        1,
				Visible:        true,
				Width:          2,
				Height:         2,
				GridSize:       100,
				Identifier:     "Tiles",
				Opacity:        1,
				TilesetDefUid:  1,
				TilesetRelPath: "inputtilemap.png",
				Tiles: []*Tile{
					&Tile{
						Position: []int{0, 0},
						Src:      []int{0, 0},
						ID:       0,
					},
					&Tile{
						Position: []int{100, 0},
						Src:      []int{0, 0},
						ID:       0,
					},
					&Tile{
						Position: []int{100, 0},
						Src:      []int{0, 0},
						ID:       0,
					},
					&Tile{
						Position: []int{100, 100},
						Src:      []int{0, 0},
						ID:       0,
					},
				},
			},
		},
	}

	// add it to existing one
	levels = append(levels, mylevel)

	// insert into primary JSON
	finaljson, _ := sjson.Set(jsonstr, "levels", levels)

	fmt.Println(finaljson)
	os.Exit(0)

	for _, level := range levels {
		if level.Identifier != "Input_1" {
			continue
		}

		for _, layer := range level.Layers {
			if layer.Identifier != "Tiles" {
				continue
			}

			for _, tile := range layer.Tiles {
				fmt.Printf("tile %d at %d,%d uses tile from %s at %d,%d\n",
					tile.ID, tile.Position[0], tile.Position[1], layer.TilesetRelPath,
					tile.Src[0], tile.Src[1],
				)
			}
		}
	}
}

func readfile(filename string) []byte {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		panic(err)
	}

	return buffer
}
