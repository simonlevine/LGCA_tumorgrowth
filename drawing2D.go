package main

import (
	"fmt"
	"image"
	"strconv"
)

// The following code was written by Simon Levine-Gottreich

//DrawMatrices takes in a slice of matrices to output slice of images that can be used to draw GIF
func DrawMatrices(matrices []Matrix2D, cellWidth int, x, y int) []image.Image {
	numGenerations := len(matrices)
	imageList := make([]image.Image, numGenerations)
	for i := range matrices {
		fmt.Println("Drawing " + strconv.Itoa(i) + "th matrix")
		imageList[i] = DrawMatrix2D(matrices[i], cellWidth, x, y)
	}
	return imageList
}

//DrawMatrix2D takes in matrix and outputs image.Image
func DrawMatrix2D(matrix Matrix2D, cellWidth int, x, y int) image.Image {
	height := len(matrix) * cellWidth
	width := len(matrix[0]) * cellWidth
	c := CreateNewCanvas(width, height)

	// declare colors
	darkGray := MakeColor(50, 50, 50)
	black := MakeColor(0, 0, 0)
	blue := MakeColor(0, 0, 255)
	red := MakeColor(255, 0, 0)
	green := MakeColor(213, 245, 227)
	yellow := MakeColor(255, 255, 0)
	white := MakeColor(255, 255, 255)

	// draw the grid lines
	c.SetStrokeColor(darkGray)
	DrawGridLines(c, cellWidth)

	// fill in colored squares
	for i := range matrix {
		for j := range matrix[i] {
			if InField2D(i, j, x, y) == true {
				if matrix[i][j].state == "C" {
					c.SetFillColor(blue)
				} else if matrix[i][j].state == "Q" {
					c.SetFillColor(yellow)
				} else if matrix[i][j].state == "h" {
					c.SetFillColor(green)

				} else if matrix[i][j].state == "N" {
					c.SetFillColor(red)

				} else if matrix[i][j].state == "wN" {
					c.SetFillColor(black)

				}
			} else {
				c.SetFillColor(white)
			}
			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}

	return c.img
}

//DrawGridLines draws gridlines
func DrawGridLines(pic Canvas, cellWidth int) {
	w, h := pic.width, pic.height
	// first, draw vertical lines
	for i := 1; i < pic.width/cellWidth; i++ {
		y := i * cellWidth
		pic.MoveTo(0.0, float64(y))
		pic.LineTo(float64(w), float64(y))
	}
	// next, draw horizontal lines
	for j := 1; j < pic.height/cellWidth; j++ {
		x := j * cellWidth
		pic.MoveTo(float64(x), 0.0)
		pic.LineTo(float64(x), float64(h))
	}
	pic.Stroke()
}
