package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
)

func main() {
	fmt.Printf("here")
	reader, err := os.Open("img.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := img.Bounds()

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			fmt.Printf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
			fmt.Printf("%v %v %v %v", r, g, b, a)
		}
	}
}
