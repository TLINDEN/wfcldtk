package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func Exists[K comparable, V any](m map[K]V, v K) bool {
	if _, ok := m[v]; ok {
		return true
	}
	return false
}

func GetTileFromSpriteSheet(img image.Image, x, y, width, height int) (image.Image, error) {
	// Create new image
	outputImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Copy pixels
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			outputImg.Set(i, j, img.At(x+i, y+j))
		}
	}

	return outputImg, nil
}

func Loadimage(filename string) (image.Image, error) {
	raw, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	img, _, err := image.Decode(raw)
	if err != nil {
		return nil, err
	}
	return img, nil

}

type Color [4]uint8

// get color value from 1 pixel position
func GetColor(img image.Image, x, y int) Color {
	r, g, b, a := img.At(x, y).RGBA()
	return Color{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// turn color value into hex string
func HexFromColor(c Color) string {
	return fmt.Sprintf("%02x%02x%02x%02x", c[0], c[1], c[2], c[3])
}

// returns the adjacency constraint id for the given tile image
// in the provided direction.
func GetConstraint(tileimage image.Image, direction Direction, count int) string {
	width := tileimage.Bounds().Max.X
	height := tileimage.Bounds().Max.Y

	wdist := width / count
	hdist := height / count

	var hash string
	points := make([]Color, count)

	for i := 0; i < count; i++ {
		switch direction {
		case North:
			points[i] = GetColor(tileimage, wdist+i*wdist, 0)
		case South:
			points[i] = GetColor(tileimage, wdist+i*wdist, height-1)
		case West:
			points[i] = GetColor(tileimage, 0, hdist+i*hdist)
		case East:
			points[i] = GetColor(tileimage, width-1, hdist+i*hdist)
		}
	}

	// Generate a hash from the colors
	hash = ""
	for _, c := range points {
		hash += HexFromColor(c)
	}

	sum := sha256.Sum256([]byte(hash))
	res := fmt.Sprintf("%x", sum)[:8]

	return res
}

// convert a color value to a byte array
func color2byte(color uint32) []byte {
	bytearray := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytearray, color)
	return bytearray
}

// calculate a hash for the whole image
func GetImageHash(tile image.Image) string {
	hash := sha256.New()

	for y := tile.Bounds().Min.Y; y < tile.Bounds().Dy(); y++ {
		for x := tile.Bounds().Min.X; x < tile.Bounds().Dx(); x++ {
			r, g, b, a := tile.At(x, y).RGBA()
			for _, color := range []uint32{r, g, b, a} {
				_, err := hash.Write(color2byte(color))
				if err != nil {
					log.Fatalf("failed to calculate image checksum: %s", err)
				}
			}
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil))[:32]
}

// check if an image is completely transparent
func ImageIsTransparent(tile image.Image) bool {
	for y := tile.Bounds().Min.Y; y < tile.Bounds().Dy(); y++ {
		for x := tile.Bounds().Min.X; x < tile.Bounds().Dx(); x++ {
			_, _, _, alpha := tile.At(x, y).RGBA()
			if alpha != 0 {
				return false
			}
		}
	}
	return true
}

func SavePNG(filename string, img image.Image) {
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}
