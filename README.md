# Wave Function Collapse 2D level generator for LDTK

## Introduction

This little tool aims to generate levels using the wave function
collapse algorithm. It takes its input from a level in a LDTK project
file.

So just create a new LDTK project, add a tileset, a tileset layer and
add a level. Give it a name and add all the tiles you wish to use for
generated levels. If you want a increase the probability, that a
particular tile appears, just add it multiple times to the level. The
tiles in this level don't need to match together.

The tool is work in progress.

## Example

This level has been generated using [this tileset](https://github.com/TLINDEN/wfcldtk/blob/main/images/inputtilemap.png)

![example-output](https://github.com/TLINDEN/wfcldtk/blob/main/images/output.png)

The commandline to generate it was:

```shell
./wfcldtk -p images/demo.ldtk -l Input_1 images/output.png -W 24 -H 12 && display images/output.png
```

## TODO
- add LDTK project file writing
- add another Populate() function to be able to pre-populate the output map using an LDTK level
- add weight to tiles in slot

## Usage

```default
This is wfcldtk, a WFC level generator for LDTK.

Usage: gfn [-vd] -p <project> -l <level> [-W <width> -H <height>] [<output image>]

Options:
-p --project <project>  Read data from LDTK file <project>
-l --level <level>      Use level <level> as example for overlap mode
-W --width <width>      Width in number of tiles (not pixel!)
-H --height <height>    Height

-d --debug    Show debugging output
-v --version  Show program version
```


## Installation

You will need the Golang toolchain  in order to build from source. GNU
Make will also help but is not strictly neccessary.

If you want to compile the tool yourself, use `git clone` to clone the
repository.   Then   execute   `go mod tidy`   to   install   all
dependencies. Then  just enter `go build` or -  if you have  GNU Make
installed - `make`.

To install after building either copy the binary or execute 
`sudo make install`. 


# Report bugs

[Please open an issue](https://github.com/TLINDEN/gfn/issues). Thanks!

# License

This work is licensed under the terms of the General Public Licens
version 3.

# Author

Copyleft (c) 2024 Thomas von Dein
