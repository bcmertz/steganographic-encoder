package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	reader, err := os.Open("img.png")
	if err != nil {
		os.Exit(1)
		log.Fatal(err)
	}
	defer reader.Close()

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	inputText := "hello world this is bennett"

	inputBinary := stringToBinary(inputText)
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	new_image := image.NewRGBA(bounds)
	// ravioli ravioli here is the formuoli: [ r - g - b - a ] - [ r - g - b - a ] - ...
	counter := 0
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixel := rgbaToBinaryPixel(int(r), int(g), int(b), int(a))
			if 4*counter+4 <= len(inputBinary) {
				bits := []string{string(inputBinary[4*counter]), string(inputBinary[4*counter+1]), string(inputBinary[4*counter+2]), string(inputBinary[4*counter+3])}
				pixel = leastBitEncoder(bits, pixel)
			}
			pixel = convertToRGBA(pixel)
			new_image.Set(x, y, color.RGBA{uint8(pixel.R), uint8(pixel.G), uint8(pixel.B), uint8(pixel.A)})
		}
	}

	// time to decode last bit by last bit!
	// throwAwayCode // here it go for that debug life
	//lastBits := getLastBits(pixels)
	//str := binaryToString([]byte(lastBits))

	output, _ := os.Create("newimage.png")
	png.Encode(output, new_image)
}

func throwAwayCode() {
	txt := "hi"
	txty := stringToBinary(txt)
	var dixels [][]Pixel
	//dix := []Pixel{Pixel{11010110, 11011100, 00110110, 00011010}, Pixel{11010110, 00111000, 10101000, 11001110}, Pixel{10101010, 11100010, 11001111, 10101010}, Pixel{10000101, 11000011, 11111110, 00000000}}
	// implementation works for every 8 bit pixel the same r=g=b=a, idk about pixel = pixel where r=/=g=/=b=/=a but this indicates a bug in the encoder
	dix := []Pixel{
		Pixel{1 * 257, 1 * 257, 2 * 257, 2 * 257},
		Pixel{2 * 257, 2 * 257, 2 * 257, 2 * 257},
		Pixel{2 * 257, 2 * 257, 2 * 257, 2 * 257},
		Pixel{2 * 257, 2 * 257, 2 * 257, 2 * 257}}

	dixels = append(dixels, dix)
	for i, pixl := range dixels[0] {
		bits := []string{string(txty[4*i]), string(txty[4*i+1]), string(txty[4*i+2]), string(txty[4*i+3])}
		//fmt.Println(1, pixl)
		pixel := rgbaToBinaryPixel(int(pixl.R), int(pixl.G), int(pixl.B), int(pixl.A)) // 257 * 256 val rgba -> 8 bit binary pixel (256)
		//fmt.Println(2, pixel)                                                          //
		new_pixel := leastBitEncoder(bits, pixel) // 8 bit -> encoded truncated int (no leading 0s)
		//fmt.Println(3, new_pixel, "\n")                                                //
		dixels[0][i] = convertToRGBA(new_pixel) // int -> 257 * 256 val rgba

		//dixels[0][i] = convertToRGBA(leastBitEncoder(bits, pixl))
	}
	//fmt.Println(dixels)
	lastBits := getLastBits(dixels)
	//fmt.Println(lastBits)
	str := binaryToString([]byte(lastBits))
	fmt.Println("str", str)
}

// WORKS
func binaryToString(s []byte) string {
	output := make([]byte, len(s)/8)
	for i := 0; i < len(output); i++ {
		val, _ := strconv.ParseInt(string(s[i*8:(i+1)*8]), 2, 64)
		output[i] = byte(val)
	}
	return string(output)
}

// WORKS
func lastBit(x int) string {
	if x%2 == 0 {
		return "0"
	} else {
		return "1"
	}
}

// WORKS
func getLastBits(pixels [][]Pixel) string {
	lastBits := ""
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			p := pixels[i][j]
			pixel := rgbaToBinaryPixel(p.R, p.G, p.B, p.A)
			lastBits += lastBit(pixel.R)
			lastBits += lastBit(pixel.G)
			lastBits += lastBit(pixel.B)
			lastBits += lastBit(pixel.A)
		}
	}
	return lastBits
}

// WORKS
func rgbaToBinaryPixel(r int, g int, b int, a int) Pixel {
	r, g, b, a = base10toBase2(r/257), base10toBase2(g/257), base10toBase2(b/257), base10toBase2(a/257)
	return Pixel{r, g, b, a}
}

// WORKS
func base10toBase2(input int) int {
	// 128,64,32,16,8,4,2,1
	//   1, 0, 0, 0,0,0,0,0
	// 128 -> 64[0] -> 32[00] -> 16[000] -> 8[0000] -> 4[00000] -> 2[000000] -> 1[0000000] -> 0[10000000]
	str := ""
	for input != 0 {
		remainder := input % 2
		str = strconv.Itoa(remainder) + str
		input = input / 2
	}
	output, _ := strconv.Atoi(str)
	return output
}

// CORRECT
type Pixel struct {
	R int
	G int
	B int
	A int
}

// AAAHHHHHHH
func leftPad(str string, length int) string {
	for len(str) < length {
		str = "0" + str
	}
	return str
}

// WORKS
func base2toBase10(input int) int {
	str := strconv.Itoa(input) // 10000000 [128]
	str = leftPad(str, 8)
	//fmt.Println("STR", str)
	accumulator := 0
	for i := len(str) - 1; i >= 0; i-- {
		char_val, _ := strconv.Atoi(string(str[i])) // = 0
		power := float64(7 - i)                     // =
		accumulator += char_val * int(math.Pow(2, power))
	}
	return accumulator
}

// WORKS
func convertToRGBA(p Pixel) Pixel {
	rgbaPixel := Pixel{base2toBase10(p.R) * 257, base2toBase10(p.G) * 257, base2toBase10(p.B) * 257, base2toBase10(p.A) * 257}
	return rgbaPixel
}

func leastBitEncoder(bin []string, pix Pixel) Pixel {
	R, _ := strconv.Atoi(bin[0])
	G, _ := strconv.Atoi(bin[1])
	B, _ := strconv.Atoi(bin[2])
	A, _ := strconv.Atoi(bin[3])
	if pix.R%2 == 0 {
		if R != 0 {
			pix.R += 1
		}
	} else {
		if R != 1 {
			pix.R -= 1
		}
	}
	if pix.G%2 == 0 {
		if G != 0 {
			pix.G += 1
		}
	} else {
		if G != 1 {
			pix.G -= 1
		}
	}
	if pix.B%2 == 0 {
		if B != 0 {
			pix.B += 1
		}
	} else {
		if B != 1 {
			pix.B -= 1
		}
	}
	if pix.A%2 == 0 {
		if A != 0 {
			pix.A += 1
		}
	} else {
		if A != 1 {
			pix.A -= 1
		}
	}

	return pix
}

func stringToBinary(str string) (output string) {
	for _, str_char := range str {
		var tmp string
		tmp = fmt.Sprintf("%s%.8b", tmp, str_char)
		output += tmp
	}
	return
}
