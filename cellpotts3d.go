package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

// The following code was written by Simon Levine-Gottreich

// Except for the lines commented //Noah Chang ----- ... //-----

// for detailed comments, refer to the 2D version

//Cell is an object with necessary information for 3D cellular automata to run
type Cell struct {

	//state is either "C", "Q", "N", "wN"
	state string

	//The location of in the matrix
	location OrderedTrio

	//The direction of the velocity
	velocityDirection OrderedTrio

	//Probability of N, A, C, Q
	pNecrosis, pProliferation, pQuiescent float64
}

//Neighborhood is an object to organize the neighbors
type Neighborhood struct { //just organizing, so using pointers.

	//Slice of pointers to adjaecent cells
	neighbors []*Cell

	//The center cell
	center *Cell
}

//Matrix is 3 dimensional slice of Cells.
type Matrix [][][]Cell

//OrderedTrio is the (x,y,z) coordinate
type OrderedTrio struct {
	x, y, z int
}

//GenerateMatrices is the 3D version of Generate2DMatrices
func GenerateMatrices(initialMatrix Matrix, numGens int, Kcc, Knn, Knc float64) []Matrix {

	matrices := make([]Matrix, numGens+1)

	matrices[0] = initialMatrix

	centerCell := GetCentralCell(matrices[0])

	matrices[0][centerCell.location.x][centerCell.location.y][centerCell.location.z].state = "C"

	for m := 1; m <= numGens; m++ {
		fmt.Println("3D Matrix Generation No." + strconv.Itoa(m))
		matrices[m] = UpdateMatrix(matrices[m-1], Kcc, Knn, Knc) //assumes moore neighborhood
	}

	return matrices
}

//GetCentralCell is the 3D version of GetCentralCell2D
func GetCentralCell(currMatrix Matrix) Cell {

	numRows := GetNumRows(currMatrix)
	numCols := GetNumCols(currMatrix)
	numAisles := GetNumAisles(currMatrix)

	centralCell := currMatrix[numRows/2][numCols/2][numAisles/2]

	return centralCell
}

//UpdateMatrix is the 3D version of Update2DMatrix
func UpdateMatrix(currMatrix Matrix, Kcc, Knn, Knc float64) Matrix {

	numRows := GetNumRows(currMatrix)
	numCols := GetNumCols(currMatrix)
	numAisles := GetNumAisles(currMatrix)

	freshMatrix := Initialize3DMatrix(numRows, numCols, numAisles) //should be zero values for states to begin with, except for seed values in center.

	freshMatrix = UpdateMatrixStates(currMatrix, freshMatrix, Kcc, Knn, Knc)

	freshMatrix = UpdateMatrixVelocities(currMatrix, freshMatrix) //BUG HERE --> index out of range.

	//now that states and velocities at each cell are reflected in freshMatrix, it's time to proliferate cells if applicable.

	freshMatrix = PushAllCells(freshMatrix)

	return freshMatrix
}

//PushAllCells is the 3D version of PushAllCells2D
func PushAllCells(currMatrix Matrix) Matrix {

	x := GetNumRows(currMatrix)
	y := GetNumCols(currMatrix)
	z := GetNumAisles(currMatrix)
	pushedMatrix := Initialize3DMatrix(x, y, z)

	for i := range currMatrix {
		for j := range currMatrix[i] {
			for k := range currMatrix[i][j] {
				pushedMatrix[i][j][k] = currMatrix[i][j][k]
			}
		}
	}

	for i := range currMatrix {
		for j := range currMatrix[i] {
			for k := range currMatrix[j][i] {

				currCell := currMatrix[i][j][k]

				toX := currCell.velocityDirection.x
				toY := currCell.velocityDirection.y
				toZ := currCell.velocityDirection.z

				if currCell.state == "C" {
					pushedMatrix[i][j][k].state = "C"
					pushedMatrix[toX][toY][toZ].state = "C"
				}

				if currCell.state == "N" {
					pushedMatrix[i][j][k].state = "h"
					pushedMatrix[toX][toY][toZ].state = "N"
				}

			}
		}
	}

	return pushedMatrix

}

//UpdateMatrixVelocities is the 3D version of UpdateMatrixVelocities2D
func UpdateMatrixVelocities(currMatrix, freshMatrix Matrix) Matrix {

	numRows := GetNumRows(currMatrix)
	numCols := GetNumCols(currMatrix)
	numAisles := GetNumAisles(currMatrix)

	for i := range currMatrix {
		for j := range currMatrix[i] {
			for k := range currMatrix[j][i] { //ranging into cell

				if InField3D(i, j, k, numRows, numCols, numAisles) == true {
					freshMatrix[i][j][k] = UpdateOneCellVelocity(currMatrix, i, j, k) // updating cell states
				}
			}
		}
	}
	return freshMatrix
}

//UpdateOneCellVelocity is the 3D version of UpdateOneCellVelocity
func UpdateOneCellVelocity(currMatrix Matrix, i, j, k int) Cell {

	currCell := currMatrix[i][j][k]

	numRows := GetNumRows(currMatrix)
	numCols := GetNumCols(currMatrix)
	numAisles := GetNumAisles(currMatrix)

	if InField3D(i, j, k, numRows, numCols, numAisles) == true {

		if currCell.state == "C" {
			//cell velocity should point to this (x,y,z).
			minCNeighborCoord := GetMinCNeighborDirection(currMatrix, i, j, k)

			//set velocity vector to point to direction of neighbor least-dense with cancer cells.
			currCell.velocityDirection = minCNeighborCoord

		} else if currCell.state == "N" {

			maxNNeighborCoord := GetMaxNNeighborDirection(currMatrix, i, j, k)
			//set velocity vector to point to neighbor with most necrosis in its neighborhood.
			currCell.velocityDirection = maxNNeighborCoord
		} else {
			currCell.velocityDirection.x = currCell.location.x
			currCell.velocityDirection.y = currCell.location.y
			currCell.velocityDirection.z = currCell.location.z
		}
	}
	return currCell
}

//GetMaxNNeighborDirection is the 3D version of GetMaxNNeighborDirection2D
func GetMaxNNeighborDirection(currMatrix Matrix, i, j, k int) OrderedTrio {

	numRows := GetNumRows(currMatrix)
	numCols := GetNumCols(currMatrix)
	numAisles := GetNumAisles(currMatrix)

	//Noah Chang------------------------------------------------------------------
	neighborhoods := make([]Neighborhood, 0)
	if InField3D(i+2, j, k, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i+2, j, k, numRows, numCols, numAisles))
	}
	if InField3D(i-2, j, k, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i-2, j, k, numRows, numCols, numAisles))
	}
	if InField3D(i, j+2, k, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i, j+2, k, numRows, numCols, numAisles))
	}
	if InField3D(i, j-2, k, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i, j-2, k, numRows, numCols, numAisles))
	}
	if InField3D(i, j, k+2, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i, j, k+2, numRows, numCols, numAisles))
	}
	if InField3D(i, j, k-2, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i, j, k+2, numRows, numCols, numAisles))
	}
	//----------------------------------------------------------------------------

	var maxNcoords OrderedTrio

	// zero the necrotic count since we want the max. Necrotic cells are chemotactic to others.
	maxCountN := 0.0

	//ranging over surrounding nhds.
	for i := range neighborhoods {

		currCountN := GetNumNecrotic(neighborhoods[i])

		//if a new MINIMUM is found, set coordinates (this is where we WANT a cancerous cell to go).
		if currCountN > maxCountN {

			maxNcoords.x = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.x
			maxNcoords.y = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.y
			maxNcoords.z = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.z

		} else if currCountN == maxCountN {

			// TIEBREAKING: we take at random the maximum N count.
			maxCountN = currCountN

			maxNcoords.x = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.x
			maxNcoords.y = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.y
			maxNcoords.z = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.z

		}

	}
	return maxNcoords
}

//GetMinCNeighborDirection is the 3D version of GetMinCNeighborDirection2D
func GetMinCNeighborDirection(currMatrix Matrix, i, j, k int) OrderedTrio {

	numRows := GetNumRows(currMatrix)
	numCols := GetNumCols(currMatrix)
	numAisles := GetNumAisles(currMatrix)

	//ranging over all neighboring cells to see which neighbor is itself surrounded by fewest cancer cells.

	var minCcoords OrderedTrio

	// make a large number out of cubed dimensions of matrix
	minCountC := float64(numRows * numCols * numAisles)

	//Noah Chang------------------------------------------------------------------
	neighborhoods := make([]Neighborhood, 0)
	if InField3D(i+2, j, k, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i+2, j, k, numRows, numCols, numAisles))
	}
	if InField3D(i-2, j, k, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i-2, j, k, numRows, numCols, numAisles))
	}
	if InField3D(i, j+2, k, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i, j+2, k, numRows, numCols, numAisles))
	}
	if InField3D(i, j-2, k, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i, j-2, k, numRows, numCols, numAisles))
	}
	if InField3D(i, j, k+2, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i, j, k+2, numRows, numCols, numAisles))
	}
	if InField3D(i, j, k-2, numRows, numCols, numAisles) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood(currMatrix, i, j, k+2, numRows, numCols, numAisles))
	}
	//----------------------------------------------------------------------------

	//ranging over surrounding nhds.
	for i := range neighborhoods {

		currCountC := GetNumCancerous(neighborhoods[i])

		if currCountC < minCountC {

			minCcoords.x = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.x
			minCcoords.y = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.y
			minCcoords.z = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.z

		} else if currCountC == minCountC {

			minCountC = currCountC

			minCcoords.x = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.x
			minCcoords.y = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.y
			minCcoords.z = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.z

		}

	}

	return minCcoords

}

//UpdateMatrixStates is the 3D version of UpdateMatrixStates2D
func UpdateMatrixStates(currMatrix, freshMatrix Matrix, Kcc, Knn, Knc float64) Matrix {

	for i := range currMatrix {
		for j := range currMatrix[i] {
			for k := range currMatrix[j][i] { //ranging into cell

				freshMatrix[i][j][k] = UpdateOneCellState(currMatrix, i, j, k, Kcc, Knn, Knc) // updating cell states

			}
		}
	}
	return freshMatrix
}

//UpdateOneCellState is the 3D version of UpdateOneCellState2D
func UpdateOneCellState(currMatrix Matrix, i, j, k int, Kcc, Knn, Knc float64) Cell {

	//Updating cell states based on previous probabilities
	// and adding current probabilities to Cell struct

	numRows := GetNumRows(currMatrix)
	numCols := GetNumCols(currMatrix)
	numAisles := GetNumAisles(currMatrix)

	currNhd := GetCurrentNeighborhood(currMatrix, i, j, k, numRows, numCols, numAisles)

	N := GetNumNecrotic(currNhd)
	C := GetNumCancerous(currNhd)

	Ep := EProliferation(Kcc, Knn, Knc, N, C)
	Eq := EQuiescence(Kcc, Knn, Knc, N, C)
	En := ENecrosis(Kcc, Knn, Knc, N, C)

	pN := ProbNecrosis(Ep, En, Eq)
	pP := ProbProliferation(Ep, En, Eq) * 10000000000.0
	pQ := ProbQuiescence(Ep, En, Eq) * 10000000.0

	var newCell Cell

	newCell = currMatrix[i][j][k] //make a copy

	currMatrix[i][j][k].pNecrosis = pN
	currMatrix[i][j][k].pQuiescent = pQ
	currMatrix[i][j][k].pProliferation = pP

	pAll := []float64{pN, pP, pQ}

	maxP := GetMaxP(pAll)

	//traceback:
	//quiescent, necrotic, proliferative are possible next states
	//new state is max of current probabilities

	if maxP == pN && C >= 1 { //cell dies only if the number of cancer cells present is greater than 1
		newCell.state = "N"
	} else if maxP == pP && C >= 1 && C+N < 5 {
		newCell.state = "C"
		//ordering this last will cause cell to default to quiescent in case of a tie.
	} else if maxP == pQ {
		newCell.state = "Q"
	}

	//now returning an identical cell, except with an updated state based upon probability of transition.
	return newCell
}

//LatticeConfigEnergy is the 3D version of LatticeConfigEnergy2D
func LatticeConfigEnergy(currMatrix Matrix, Kcc, Knn, Knc float64) float64 {

	latticeConfigEnergy := 0.0 //sum of energy over all neighborhoods

	numRows := GetNumRows(currMatrix)
	numCols := GetNumCols(currMatrix)
	numAisles := GetNumAisles(currMatrix) //getting dimensions.

	for i := range currMatrix {
		for j := range currMatrix[i] {
			for k := range currMatrix[i][j] {

				currNeighborhood := GetCurrentNeighborhood(currMatrix, i, j, k, numRows, numCols, numAisles)

				N := GetNumNecrotic(currNeighborhood)
				C := GetNumCancerous(currNeighborhood)

				currNhdEnergy := NeighborhoodConfigEnergy(currNeighborhood, Kcc, Knn, Knc, N, C)

				latticeConfigEnergy += currNhdEnergy

			}
		}
	}

	return latticeConfigEnergy
}

//Initialize3DMatrix sets up a new 3D board
func Initialize3DMatrix(x, y, z int) Matrix {

	matrix := make(Matrix, x)

	for i := range matrix {
		matrix[i] = make([][]Cell, y)

		for j := range matrix[i] {

			matrix[i][j] = make([]Cell, z)
		}
	}

	for l := range matrix { //looping through all cells and replacing nil vals.
		for m := range matrix[l] {
			for n := range matrix[l][m] {

				matrix[l][m][n].state = "h" //healthy normal cells (boundary cases).

				var cellLocation OrderedTrio
				cellLocation.x = l
				cellLocation.y = m
				cellLocation.z = n

				matrix[l][m][n].location = cellLocation

			}
		}
	}

	return matrix
}

//GetCurrentNeighborhood is the 3D version of GetCurrentNeighborhood2D
func GetCurrentNeighborhood(currMatrix Matrix, x, y, z int, numRows, numCols, numAisles int) Neighborhood { //m1

	currCenterCell := currMatrix[x][y][z]

	var currNhd Neighborhood //establishing a slice of cells that represent VonNeumann 3d Neighborhood

	currNhd.center = &currCenterCell

	if InField3D(x, y, z, numRows, numCols, numAisles) == true { //discount center cell and make sure we are in the field.

		neighborSlice := []*Cell{

			&currMatrix[x-1][y][z],
			&currMatrix[x+1][y][z],
			&currMatrix[x][y-1][z],
			&currMatrix[x][y+1][z],
			&currMatrix[x][y][z-1],
			&currMatrix[x][y][z+1],
		}

		for _, nbr := range neighborSlice {

			currNhd.neighbors = append(currNhd.neighbors, nbr)

		}
	}

	return currNhd //von-neumann
}

//InField3D checks if the given x,y,z, coordinate is in field
func InField3D(x, y, z int, numRows, numCols, numAisles int) bool {

	// since we check neighborhoods of neighbors of a given cell, the border case is TWO cells in magnitude
	if (x-5 < 0 || x+5 > numRows) || (y-5 < 0 || y+5 > numCols) || (z-5 < 0 || z+5 > numAisles) {
		return false //out of matrix field
	}
	return true //in the matrix field.
}

//GetNumRows gets the number of rows
func GetNumRows(matrix Matrix) int {
	rows := len(matrix)
	return rows
}

//GetNumCols gets the number of columns
func GetNumCols(matrix Matrix) int {
	cols := len(matrix[0])
	return cols
}

//GetNumAisles gets the number of aisles
func GetNumAisles(matrix Matrix) int {
	aisles := len(matrix[0][0])
	return aisles
}

//AssertCuboidMatrix checks if the Matrix has same lengths of rows, cols, and aisles
func AssertCuboidMatrix(currMatrix Matrix) {
	if len(currMatrix) == 0 {
		fmt.Println("Game board has no rows.")
		os.Exit(2)
	}

	numCols := len(currMatrix[0])
	for i := 1; i < len(currMatrix); i++ {
		if len(currMatrix[i]) != numCols {
			fmt.Println("Board isn't cuboid.")
			os.Exit(1)
		}
	}
}

//GetNumCancerous is the 3D version of GetNumCancerous2D
func GetNumCancerous(nhd Neighborhood) float64 {

	numC := 0.0 //initialize to zero

	for i := range nhd.neighbors {
		if nhd.neighbors[i].state == "C" || nhd.neighbors[i].state == "Q" {
			numC++
		}
	}
	//now, center cell imputed
	if nhd.center.state == "C" || nhd.center.state == "Q" {
		numC++
	}

	return numC
}

//GetNumNecrotic is the 3D version of GetNumNecrotic2D
func GetNumNecrotic(nhd Neighborhood) float64 {

	numN := 0.0 //initialize to zero

	for i := range nhd.neighbors {
		if nhd.neighbors[i].state == "N" {
			numN++
		}
	}
	if nhd.center.state == "N" {
		numN++
	}

	return numN

}
