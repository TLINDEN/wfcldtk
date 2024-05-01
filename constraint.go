package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
)

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
