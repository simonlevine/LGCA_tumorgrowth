package main

import (
	"fmt"
	"image"
	"math/rand"
	"strconv"
)

//These codes are written by Noah Chang

//Metastasis2D takes in a current Matrix, current metasized cell count, and the ruptured vessel seed type to return the cumulative number of cells metastasized.
func Metastasis2D(currMatrix Matrix2D, metaBoard [][]bool, metaCount [3]int) [3]int {

	numRows := GetNumRows2D(currMatrix)
	numCols := GetNumCols2D(currMatrix)

	for i := range currMatrix {
		for j := range currMatrix {
			//If the cell is cancerous
			if currMatrix[i][j].state == "C" {
				//And there is a ruptured vessel at the same coordinate
				if IsVascular(metaBoard, currMatrix, i, j) == true {

					nhd := GetCurrentNeighborhood2D(currMatrix, i, j, numRows, numCols)

					//case for a single cancer cell
					if GetNumCancerous2D(nhd) == 0 {
						//Does it survive inside the blood vessel?
						if SurvivalCheck("single") == true {
							//If it does, it extravastates into one of three destinations
							metaCount = Extravastate(metaCount)
						}
					}

					//case for a cluster of cancer cells
					if GetNumCancerous2D(nhd) > 0 {
						//Does it survive inside the blood vessel?
						if SurvivalCheck("cluster") == true {
							//If it does, it extravastates into one of three destinations
							metaCount = Extravastate(metaCount)
						}
					}

				}
			}
		}
	}

	return metaCount
}

//GenerateMetastasisBoard2D generates the coordinates of the ruptured vessels
func GenerateMetastasisBoard2D(currMatrix Matrix2D) [][]bool {
	numCols := GetNumCols2D(currMatrix)
	numRows := GetNumRows2D(currMatrix)

	metaBoard := make([][]bool, numRows)

	for i := range currMatrix {
		metaBoard[i] = make([]bool, numCols)
	}

	for i := range metaBoard {
		for j := range metaBoard[i] {
			metaBoard[i][j] = false
		}
	}

	return metaBoard
}

//SeedMetastasisBoard2D seeds ruptured vascular in either single random coordinate or four equidistance coordinates on the board
func SeedMetastasisBoard2D(metaBoard [][]bool, seedType string) [][]bool {

	length := len(metaBoard)
	if seedType == "random" {
		metaBoard[rand.Intn(length)][rand.Intn(length)] = true
	} else if seedType == "set" {
		metaBoard[length/4][length/2] = true
		metaBoard[length/2][length/4] = true
		metaBoard[length*3/4][length/2] = true
		metaBoard[length/2][length*3/4] = true
	} else {
		panic("Seed type has to be either random or set")
	}
	return metaBoard
}

//IsVascular checks if the given coordinate is cancerous as well as ruptured vessel
func IsVascular(metaBoard [][]bool, currMatrix Matrix2D, i, j int) bool {
	vascular := false

	if metaBoard[i][j] == true {
		if currMatrix[i][j].state == "C" {
			vascular = true
		}
	}

	return vascular
}

//SurvivalCheck calculates the probability of survival of either single cell or cluster of cells according to probability from a literature
func SurvivalCheck(nbhState string) bool {
	survived := false

	if nbhState == "single" {
		prob := rand.Intn(10000)
		if prob <= 5 {
			survived = true
			fmt.Println("A single cell has survived! Extravastating...")
		}
	}

	if nbhState == "cluster" {
		prob := rand.Intn(10000)
		if prob <= 250 {
			survived = true
			fmt.Println("A cluster of cells has survived! Extravastating...")
		}
	}

	return survived

}

//Extravastate simulates the extravastation of breast cancer cell/cells into either bone, lungs, or liver.
func Extravastate(metaCount [3]int) [3]int {
	prob := rand.Intn(10000)

	//metaCount[0]=Bones, metaCount[1]=lungs, metaCount[2]=liver

	//Bones
	if prob < 5461 {
		metaCount[0]++
	}

	//Lungs
	if prob > 5461 && prob < 5461+2553 {
		metaCount[1]++
	}

	//Liver
	if prob > 5461+2553 && prob < 9999 {
		metaCount[2]++
	}

	return metaCount

}

//Generate2DMatricesMetastasis expands on the Generate2DMatrices function and adds a metastasis part
func Generate2DMatricesMetastasis(numGens int, x, y int, Kcc, Knn, Knc float64, seedType string) ([]Matrix2D, [][3]int) {

	matrices := make([]Matrix2D, numGens+1)
	matrices[0] = Initialize2DMatrix(x, y)
	centerCell := GetCentralCell2D(matrices[0])

	cellNhd := GetCurrentNeighborhood2D(matrices[0], centerCell.location.x, centerCell.location.y, x, y)
	for n := range cellNhd.neighbors {
		matrices[0][cellNhd.neighbors[n].location.x][cellNhd.neighbors[n].location.y].state = "C"
		i := cellNhd.neighbors[n].location.x
		j := cellNhd.neighbors[n].location.y

		neighborNhd := GetCurrentNeighborhood2D(matrices[0], i, j, x, y)

		for m := range neighborNhd.neighbors {
			matrices[0][neighborNhd.neighbors[m].location.x][neighborNhd.neighbors[m].location.y].state = "C"

		}
	}
	matrices[0][centerCell.location.x][centerCell.location.y].state = "C"

	//metastasis edited code -----------------------------------------------------
	metaSlice := make([][3]int, 0)
	firstGenMeta := [3]int{0, 0, 0}
	metaSlice = append(metaSlice, firstGenMeta)

	metaBoard := GenerateMetastasisBoard2D(matrices[0])
	metaBoard = SeedMetastasisBoard2D(metaBoard, seedType)

	for m := 1; m <= numGens; m++ {
		fmt.Println("Updating " + strconv.Itoa(m) + "th generation...")
		matrices[m] = Update2DMatrix(matrices[m-1], Kcc, Knn, Knc)

		nextMetaCount := Metastasis2D(matrices[m], metaBoard, metaSlice[m-1])
		metaSlice = append(metaSlice, nextMetaCount)
	}

	return matrices, metaSlice
	//----------------------------------------------------------------------------
}

//DrawMetastasisBoard draws outputs an image.Image of a metastasis board
func DrawMetastasisBoard(matrix [][]bool, cellWidth int, x, y int) image.Image {
	height := len(matrix) * cellWidth
	width := len(matrix[0]) * cellWidth
	c := CreateNewCanvas(width, height)

	// declare colors
	red := MakeColor(255, 0, 0)
	white := MakeColor(255, 255, 255)

	// fill in colored squares
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] == true {
				c.SetFillColor(red)
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
