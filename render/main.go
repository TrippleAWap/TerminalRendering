package render

import (
	"fmt"
	"github.com/nfnt/resize"
	"golang.org/x/term"
	"image"
	"os"
	"strings"
)

var lastSize = 0

func ResetTerminal() {
	fmt.Print("\033[H\033[3J")
	// the code below is the actual code to clear the terminal screen
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Printf("Failed to get terminal size: %s\n", err)
		os.Exit(1)
	}
	if lastSize != width*height {
		buf := make([]byte, width*height)
		for i := 0; i < len(buf); i++ {
			buf[i] = ' '
		}
		_, err = os.Stdout.Write(buf)
		if err != nil {
			fmt.Printf("Failed to write terminal size: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("\033[0;0H")
		lastSize = width * height
	}
}

func RenderImage(img image.Image) string {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	aspect := float64(width) / float64(height)
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	output := strings.Builder{}
	if err != nil {
		fmt.Printf("Failed to get terminal size: %s\n", err)
		os.Exit(1)
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
