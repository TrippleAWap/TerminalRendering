package main

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"golang.org/x/term"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <image_file>\n", strings.Split(os.Args[0], "\\")[len(strings.Split(os.Args[0], "\\"))-1])
		os.Exit(1)
	}

	imagePath := os.Args[1]
	// Open the image file
	file, err := os.ReadFile(imagePath)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}
	for {
		str := RenderImage(img)
		ResetTerminal()
		_, err = os.Stdout.WriteString(str)
		if err != nil {
			panic(err)
		}
	}
}

func ResetTerminal() {
	fmt.Print("\033[H\033[3J")
	// the code below is the actual code to clear the terminal screen
	//width, height, err := term.GetSize(int(os.Stdout.Fd()))
	//if err != nil {
	//	panic(err)
	//}
	//buf := make([]byte, width*height)
	//for i := 0; i < len(buf); i++ {
	//	buf[i] = ' '
	//}
	//os.Stdout.Write(buf)
	fmt.Printf("\033[0;0H")
}

func RenderImage(img image.Image) string {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	aspect := float64(width) / float64(height)
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	output := strings.Builder{}
	if err != nil {
		panic(err)
	}
	termHeight -= 3
	if height > termHeight {
		height = termHeight
		width = int(float64(height) * aspect)
	}
	if width > termWidth {
		width = termWidth
		height = int(float64(width) / aspect)
	}
	img = resize.Resize(uint(width), uint(height), img, resize.Bilinear)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x >= img.Bounds().Dx() || y >= img.Bounds().Dy() {
				continue
			}
			r, g, b, _ := img.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			output.WriteString(fmt.Sprintf("\033[48;2;%d;%d;%dm  \033[0m", r, g, b))
		}
		output.WriteString("\n\033[0m")
	}
	output.WriteString("\033[0m")
	return output.String()
}
