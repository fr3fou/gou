package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"time"

	"github.com/JamesHovious/w32"
	"github.com/kbinani/screenshot"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

// Pou Colors

var colors = []color.Color{
	color.RGBA{99, 199, 255, 255},
	color.RGBA{247, 239, 57, 255},
	color.RGBA{189, 117, 255, 255},
	color.RGBA{255, 130, 33, 255},
	color.RGBA{66, 243, 49, 255},
	color.RGBA{255, 130, 181, 255},
	color.RGBA{140, 138, 140, 255},
	// color.RGBA{255, 255, 255, 255},
}

var down = w32.INPUT{
	Type: 0,
	Mi: w32.MOUSEINPUT{
		DwFlags: w32.MOUSEEVENTF_LEFTDOWN,
	},
}

var up = w32.INPUT{
	Type: 0,
	Mi: w32.MOUSEINPUT{
		DwFlags: w32.MOUSEEVENTF_LEFTUP,
	},
}

var run = false

func main() {
	levelsToPass, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	kbC, err := RegKbHook()
	if err != nil {
		panic(err)
	}

	fmt.Println(" Move Cursor to the TOP LEFT corner of the playground then press '1' to set TOP LEFT cords! ")
	fmt.Println()
	fmt.Println(" Move Cursor to the BOTTOM RIGHT corner of the playground then press '2' to set BOTTOM RIGHT cords! ")
	fmt.Println()
	fmt.Println(" Press '3' to start the bot! ")
	fmt.Println(" Press '4' to stop the bot! ")
	type cords struct {
		x0 int
		y0 int
		x1 int
		y1 int
	}

	c := &cords{}

	for ev := range kbC {
		if ev.Message == types.WM_KEYUP && ev.VKCode == types.VK_1 {
			x0, y0, ok := w32.GetCursorPos()
			if !ok {
				panic("Could not get Cursor Pos with win32!")
			}
			c.x0 = x0
			c.y0 = y0

			fmt.Printf("Setting x0,y0 to %+d,%+d\n", x0, y0)
		}

		if ev.Message == types.WM_KEYUP && ev.VKCode == types.VK_2 {
			x1, y1, ok := w32.GetCursorPos()
			if !ok {
				panic("Could not get Cursor Pos with win32!")
			}
			c.x1 = x1
			c.y1 = y1

			fmt.Printf("Setting x1,y1 to %+d,%+d\n", x1, y1)
		}
		if ev.Message == types.WM_KEYUP && ev.VKCode == types.VK_3 {
			Play(levelsToPass, c.x0, c.y0, c.x1, c.y1)
		}
	}

}

func MoveClick(x, y int, delay time.Duration) {
	w32.SetCursorPos(x, y)

	err := w32.SendInput([]w32.INPUT{down})
	if err != nil {
		panic(err)
	}

	time.Sleep(delay)

	err = w32.SendInput([]w32.INPUT{up})
	if err != nil {
		panic(err)
	}
}

func RegKbHook() (chan types.KeyboardEvent, error) {
	keyboardChan := make(chan types.KeyboardEvent, 11200)

	err := keyboard.Install(nil, keyboardChan)
	if err != nil {
		return nil, err
	}

	defer func() {
		fmt.Println("here")
		_ = keyboard.Uninstall()
	}()

	return keyboardChan, nil
}

func Play(levelsToPass, x0, y0, x1, y1 int) {
	bounds := image.Rect(x0, y0, x1, y1)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	screenshotBounds := img.Bounds()
	width := screenshotBounds.Max.X
	height := screenshotBounds.Max.Y

	for i := 1; i <= levelsToPass; i++ {
		start := time.Now()

		fmt.Printf("Playing level %d/%d!\n", i, levelsToPass)

		var pouColors = make(map[color.Color]int)
		var pouColor color.Color = color.Transparent

		// This loop gets all pou colors from current level
		for y := screenshotBounds.Min.Y; y < height; y += 20 {
			for x := screenshotBounds.Min.X; x < width; x += 20 {
				pix := img.At(x, y)

				for _, c := range colors {
					if c == pix {
						pouColors[c]++
					}
				}

			}
		}

		// This loop sets pouColor (the color to look for - less is better)
		min := int(^uint(0) >> 1)
		for k, v := range pouColors {
			if v < min {
				min = v
				pouColor = k
			}
		}

		for y := screenshotBounds.Min.Y; y < height; y += 20 {
			for x := screenshotBounds.Min.X; x < width; x += 20 {

				if !run {
					return
				}
				pix := img.At(x, y)

				if pouColor == pix {
					MoveClick(x0+x, y0+y, time.Millisecond*60)

					img, err = screenshot.CaptureRect(bounds)
					if err != nil {
						panic(err)
					}
				}
			}
		}

		elapsed := time.Since(start)
		fmt.Println()
		fmt.Printf("Level %d passed in %s!\n", i, elapsed)
	}
}
