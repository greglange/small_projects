package main

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: mkmaze <image|video>")
		return
	}

	runType := os.Args[1]
	if !(runType == "image" || runType == "video") {
		fmt.Fprintln(os.Stderr, "Usage: mkmaze <image|video>")
		return
	}

	var err error
	if runType == "image" {
		err = mainImage()
	} else if runType == "video" {
		err = mainVideo()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func mainImage() error {
	return makeMaze(false)
}

func mainVideo() error {
	return makeMaze(true)
}

func makeMaze(video bool) error {
	cellSize := 50.0
	imageWidth, imageHeight := 1920, 1080
	cellsWidth, cellsHeight := imageWidth/int(cellSize)-2, imageHeight/int(cellSize)-2
	xBorder := (imageWidth - (cellsWidth-1)*int(cellSize)) / 2
	yBorder := (imageHeight - (cellsHeight-1)*int(cellSize)) / 2
	originX, originY := cellsWidth-1, cellsHeight-1

	const (
		none = iota
		up
		right
		down
		left
	)

	cells := make([][]int, cellsWidth)
	for i := 0; i < cellsWidth; i++ {
		cells[i] = make([]int, cellsHeight)
	}

	for x := 0; x < cellsWidth; x++ {
		for y := 0; y < cellsHeight; y++ {
			if y < cellsHeight-1 {
				cells[x][y] = down
			} else {
				cells[x][y] = right
			}
			if x == cellsWidth-1 && y == cellsHeight-1 {
				cells[x][y] = none
			}
		}
	}

	makeImage := func() *gg.Context {
		dc := gg.NewContext(imageWidth, imageHeight)
		dc.DrawRectangle(0, 0, float64(imageWidth-1), float64(imageHeight-1))
		dc.SetRGB255(255, 255, 255)
		dc.Fill()
		for x := 0; x < cellsWidth; x++ {
			for y := 0; y < cellsHeight; y++ {
				xx := float64(x*int(cellSize) + xBorder)
				yy := float64(y*int(cellSize) + yBorder)
				dc.DrawCircle(xx, yy, 8)
				if x == originX && y == originY {
					dc.SetRGB255(255, 0, 0)
				} else {
					dc.SetRGB255(0, 0, 0)
				}
				dc.Fill()
				if cells[x][y] == none {
					continue
				}
				switch cells[x][y] {
				case up:
					dc.SetLineWidth(4)
					dc.DrawLine(xx, yy, xx, yy-cellSize)
					dc.Stroke()
				case right:
					dc.SetLineWidth(4)
					dc.DrawLine(xx, yy, xx+cellSize, yy)
					dc.Stroke()
				case down:
					dc.SetLineWidth(4)
					dc.DrawLine(xx, yy, xx, yy+cellSize)
					dc.Stroke()
				case left:
					dc.SetLineWidth(4)
					dc.DrawLine(xx, yy, xx-cellSize, yy)
					dc.Stroke()
				}
			}
		}
		return dc
	}

	for i := 0; i < cellsWidth*cellsHeight*10; i++ {
		if video {
			dc := makeImage()
			png.Encode(os.Stdout, dc.Image())
		}
		var move int
		newOriginX, newOriginY := originX, originY
		if 0 == rand.Intn(2) {
			// horizontal move
			if originX == 0 {
				move = right
				newOriginX += 1
			} else if originX == cellsWidth-1 {
				move = left
				newOriginX -= 1
			} else {
				if 0 == rand.Intn(2) {
					move = right
					newOriginX += 1
				} else {
					move = left
					newOriginX -= 1
				}
			}
		} else {
			// vertical move
			if originY == 0 {
				move = down
				newOriginY += 1
			} else if originY == cellsHeight-1 {
				move = up
				newOriginY -= 1
			} else {
				if 0 == rand.Intn(2) {
					move = down
					newOriginY += 1
				} else {
					move = up
					newOriginY -= 1
				}
			}
		}
		cells[originX][originY] = move
		originX, originY = newOriginX, newOriginY
		cells[originX][originY] = none
	}
	dc := makeImage()
	png.Encode(os.Stdout, dc.Image())
	return nil
}
