/*
Copyright © 2024 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

type Template struct {
	Name string
	Tmpl string
}

const (
	VERSION string = "0.0.1"
	Usage   string = `This is wfcldtk, a WFC level generator for LDTK.

Usage: gfn [-vd] -p <project> -l <level> [-W <width> -H <height>] [<output image>]

Options:
-p --project <project>  Read data from LDTK file <project>
-l --level <level>      Use level <level> as example for overlap mode
-W --width <width>      Width in number of tiles (not pixel!)
-H --height <height>    Height

-d --debug    Show debugging output
-v --version  Show program version

`
	DefaultWidth       int = 4
	DefaultHeight      int = 4
	DefaultCheckpoints int = 5
)

type Config struct {
	Showversion bool   `koanf:"version"` // -v
	Debug       bool   `koanf:"debug"`   // -d
	Project     string `koanf:"project"`
	Level       string `koanf:"level"`
	Height      int    `koanf:"height"`
	Width       int    `koanf:"width"`
	Outputimage string // arg 1 just used for debugging currently
	Checkpoints int    `koanf:"checkpoints"`
}

func InitConfig(output io.Writer) (*Config, error) {
	var kloader = koanf.New(".")

	// Load default values using the confmap provider.
	if err := kloader.Load(confmap.Provider(map[string]interface{}{
		"width":       DefaultWidth,
		"height":      DefaultHeight,
		"checkpoints": DefaultCheckpoints,
	}, "."), nil); err != nil {
		return nil, fmt.Errorf("failed to load default values into koanf: %w", err)
	}

	// setup custom usage
	flagset := flag.NewFlagSet("config", flag.ContinueOnError)
	flagset.Usage = func() {
		fmt.Fprintln(output, Usage)
		os.Exit(0)
	}

	// parse commandline flags
	flagset.BoolP("version", "v", false, "show program version")
	flagset.BoolP("debug", "d", false, "enable debug output")
	flagset.IntP("width", "W", 0, "output width")
	flagset.IntP("height", "H", 0, "output height")
	flagset.StringP("project", "p", "", "LDTK project file")
	flagset.StringP("level", "l", "", "LDTK level")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse program arguments: %w", err)
	}

	// command line setup
	if err := kloader.Load(posflag.Provider(flagset, ".", kloader), nil); err != nil {
		return nil, fmt.Errorf("error loading flags: %w", err)
	}

	// fetch values
	conf := &Config{}
	if err := kloader.Unmarshal("", &conf); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	// arg is the output file
	if len(flagset.Args()) > 0 {
		conf.Outputimage = flagset.Args()[0]
	}

	return conf, nil
}
