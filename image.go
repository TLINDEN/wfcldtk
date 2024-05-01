package main

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func SavePNG(filename string, img image.Image) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
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
