package main

import (
	"fmt"
	_ "image/png"
	"io"
	"log"
	"os"
)

func Die(err error) int {
	log.Fatal("Error: ", err.Error())

	return 1
}

func main() {
	os.Exit(Main(os.Stdout))
}

func TMain() int {
	return Main(os.Stdout)
}

func Main(output io.Writer) int {
	conf, err := InitConfig(output)
	if err != nil {
		return Die(err)
	}

	if conf.Debug {
		DEBUG = true
	}

	if conf.Project == "" || conf.Level == "" {
		Die(fmt.Errorf("mandatory parameters -p and -l missing"))
	}

	wave, err := NewWaveFromProject(conf.Project, conf.Level, conf.Width, conf.Height, conf.Checkpoints)
	if err != nil {
		log.Fatal(err)
	}

	if conf.Debug {
		fmt.Println("Superposition:")
		for _, tile := range wave.Superposition {
			fmt.Println(tile.Dump())
		}
	}

	err = wave.Collapse(100)
	if err != nil {
		Die(err)
	}

	if conf.Outputimage != "" {
		err = wave.Export(conf.Outputimage)
		if err != nil {
			log.Fatalf("failed to render: %s", err)
		}
	}

	wave.OutputTilemap.Printstats()
	fmt.Println("ok")

	return 0
}

/*

   older variant, still implemented:

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
