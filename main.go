package main

import (
	_ "image/png"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 5 {
		panic("Usage: app image width height cellsize")
	}

	width, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("failed to parse width: %s\n", err)
	}

	height, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("failed to parse height: %s\n", err)
	}

	cellsize, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatalf("failed to parse cellsize: %s\n", err)
	}

	tileset, err := Loadimage(os.Args[1])
	if err != nil {
		log.Fatalf("failed to load image: %s\n", err)
	}

	wave := NewWave(tileset, width, height, cellsize, 5)
	wave.Output.Dump()
}
