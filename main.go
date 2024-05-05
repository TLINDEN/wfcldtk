package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 6 {
		panic("Usage: ./wfcldtk ldtk-project level width height debug")
	}

	project := os.Args[1]
	level := os.Args[2]
	outputfile := os.Args[3]

	width, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatalf("failed to parse width: %s\n", err)
	}

	height, err := strconv.Atoi(os.Args[5])
	if err != nil {
		log.Fatalf("failed to parse height: %s\n", err)
	}

	if len(os.Args) > 6 {
		DEBUG = true
	}

	wave, err := NewWaveFromProject(project, level, width, height, 5)
	if err != nil {
		log.Fatal(err)
	}

	for _, tile := range wave.Superposition {
		fmt.Println(tile.Dump())
	}

	wave.Collapse(100) // ignore err for now

	err = wave.Export(outputfile)
	if err != nil {
		log.Fatalf("failed to render: %s", err)
	}

	wave.OutputTilemap.Printstats()
	fmt.Println("ok")

	// width, err := strconv.Atoi(os.Args[3])
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse width: %s\n", err)
	// }

	// height, err := strconv.Atoi(os.Args[4])
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse height: %s\n", err)
	// }

	// cellsize, err := strconv.Atoi(os.Args[5])
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse cellsize: %s\n", err)
	// }

	/*
		tileset, err := Loadimage(os.Args[1])
		if err != nil {
			return nil, fmt.Errorf("failed to load image: %s\n", err)
		}

		wave := NewWave(tileset, "", "", width, height, cellsize, 5)
		wave.Collapse(100) // ignore err for now
		wave.OutputTilemap.Dump()

		err = wave.Export(os.Args[5])
		if err != nil {
			return nil, fmt.Errorf("failed to render: %s", err)
		}

		wave.OutputTilemap.Printstats()
		fmt.Println("ok")
	*/
}
