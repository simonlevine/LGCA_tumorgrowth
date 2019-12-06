package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//Read the readme.pdf for operation

// The following code was written by Simon Levine-Gottreich

// Except for the lines commented //Noah Chang ----- ... //-----

//Cell2D is an object with necessary information for cellular automata to run
type Cell2D struct {

	//state is either "C", "Q", "N", "wN"
	state string

	//The location of in the matrix
	location OrderedPair

	//The direction of the velocity
	velocityDirection OrderedPair

	//Probability of N, A, C, Q
	pNecrosis, pProliferation, pQuiescent float64
}

//Matrix2D is a 2 dimensional slice of Cells.
type Matrix2D [][]Cell2D

//Neighborhood2D is an object to organize the neighbors
type Neighborhood2D struct {

	//Slice of pointers to adjaecent cells
	neighbors []*Cell2D

	//The center cell
	center *Cell2D
}

//OrderedPair is a (x,y) coordinate
type OrderedPair struct {
	x, y int
}

//Noah Chang--------------------------------------------------------------------
func main() {

	//seeding PRNG
	rand.Seed(time.Now().UTC().UnixNano())

	//2D Cellular automata
	if os.Args[1] == "2D" {

		fmt.Println("Cellular Automata/ Cell2D Potts/ Lattice Gas model in 2 dimensions.")

		//Establishing Boltzmann Factor constants: K_xy is a coulping constant between cell types x and y.
		//Note that similar cells have higher coupling constants to better emulate in-vivo interactions.

		//numGens - can be very high number
		numGens, _ := strconv.Atoi(os.Args[2])
		//Kcc = 3.0 recommended, per literature
		KccINT, _ := strconv.Atoi(os.Args[3])
		Kcc := float64(KccINT)
		//Knn = 3.0 recommended, per literature
		KnnINT, _ := strconv.Atoi(os.Args[4])
		Knn := float64(KnnINT)
		//Knc = 1.0 recommended, per literature ; similar cells have stronger adhesion
		KncINT, _ := strconv.Atoi(os.Args[5])
		Knc := float64(KncINT)

		fmt.Println("***************************")

		//GIF cellWidth
		cellWidth := 1

		//boardsize corresponding to spatial size of breast cancer
		x := 201
		y := 201

		//Simulation without metastasis
		if os.Args[6] == "no" {

			fmt.Println("Playing automata....")

			timepoints := Generate2DMatrices(numGens, x, y, Kcc, Knn, Knc)

			// produce animated GIF corresponding to automaton

			imglist := DrawMatrices(timepoints, cellWidth, x, y)

			outputFile := "growth"

			ImagesToGIF(imglist, outputFile)

			//Outputting CSV files for R input
			OutputFile2DinCSV(timepoints)

		}

		//Simulation with Metastasis
		if os.Args[6] == "yes" {
			fmt.Println("Playing automata with metastasis....")

			seedType := os.Args[7]

			timepoints, metaSlice := Generate2DMatricesMetastasis(numGens, x, y, Kcc, Knn, Knc, seedType)

			imglist := DrawMatrices(timepoints, cellWidth, x, y)

			outputFile := "growth"

			ImagesToGIF(imglist, outputFile)

			//Outputting CSV files for R input
			OutputFile2DinCSV(timepoints)

			//Outputting a CSV file for counting the number of cells metastasized
			OutputFileMetastasisInCSV(metaSlice)

			//Code used to draw "set" metaBoard---------------------------------------
			// metaBoard := GenerateMetastasisBoard2D(timepoints[0])
			// metaBoard = SeedMetastasisBoard2D(metaBoard, seedType)
			//
			// img := DrawMetastasisBoard(metaBoard, 1, 201, 201)
			// imgSlice := make([]image.Image, 0)
			// imgSlice = append(imgSlice, img)
			//
			// ImagesToGIF(imgSlice, "metaBoard")
			//------------------------------------------------------------------------
		}
	}

	//3D Cellular Automata
	if os.Args[1] == "3D" {

		//numGens lower than 33 recommended
		numGens, _ := strconv.Atoi(os.Args[2])
		//Kcc = 3.0 recommended, per literature
		KccINT, _ := strconv.Atoi(os.Args[3])
		Kcc := float64(KccINT)
		//Knn = 3.0 recommended, per literature
		KnnINT, _ := strconv.Atoi(os.Args[4])
		Knn := float64(KnnINT)
		//Knc = 1.0 recommended, per literature
		KncINT, _ := strconv.Atoi(os.Args[5])
		Knc := float64(KncINT)
		//Running...
		timepoints := GenerateMatrices(Initialize3DMatrix(100, 100, 100), numGens, Kcc, Knn, Knc)
		//Generating CSV for R input
		OutputFile3DinCSV(timepoints)
	}

	//2D Gif generation after R ggplot2
	if os.Args[1] == "gif2D" {
		fmt.Println("2D GIF generation")
		dir := GetNewFolderDir("outputcsv2D")
		imglist := ReadPNGs(dir)
		ImagesToGIF(imglist, "ggplot")
	}

	//3D Gif generation after R plot3D
	if os.Args[1] == "gif3D" {
		fmt.Println("3D GIF generation")
		dir := GetNewFolderDir("outputcsv3D")
		imglist := ReadPNGs(dir)
		ImagesToGIF(imglist, "ggplot3D")
	}
}

//------------------------------------------------------------------------------

//GetCentralCell2D takes in a matrix board and returns the cell at the middle of the board (2D).
func GetCentralCell2D(currMatrix Matrix2D) Cell2D { //gets cell at center of matrix. Will be used for seeding.

	numRows := GetNumRows2D(currMatrix)
	numCols := GetNumCols2D(currMatrix)

	//outputted x,y,z is center cell's row/col/aisle index of matrix

	centralCell := currMatrix[numRows/2][numCols/2] //getting to center of matrix...

	return centralCell //returns COPY of central cell. that's ok though.
}

// Generate2DMatrices is the main function of this model, this function generates numGens number of matrices for plotting according to the Lattice Gas Cellular
// Automata model. X,y are board dimensions, and Ks are coupling constants.
func Generate2DMatrices(numGens int, x, y int, Kcc, Knn, Knc float64) []Matrix2D {

	//creating slice of number of desired matrices
	matrices := make([]Matrix2D, numGens+1)

	//first matrix will be initialized...
	matrices[0] = Initialize2DMatrix(x, y)

	//seeding with cancerous cells at center.
	centerCell := GetCentralCell2D(matrices[0])

	cellNhd := GetCurrentNeighborhood2D(matrices[0], centerCell.location.x, centerCell.location.y, x, y)

	for n := range cellNhd.neighbors { //of all neighbors to given cell...

		matrices[0][cellNhd.neighbors[n].location.x][cellNhd.neighbors[n].location.y].state = "C"

		i := cellNhd.neighbors[n].location.x
		j := cellNhd.neighbors[n].location.y

		//...get the neighborhood of that cell
		neighborNhd := GetCurrentNeighborhood2D(matrices[0], i, j, x, y)

		//meta step: ranging over neighborhood of that original cell's neighbors.
		for m := range neighborNhd.neighbors {
			matrices[0][neighborNhd.neighbors[m].location.x][neighborNhd.neighbors[m].location.y].state = "C"

		}
	}

	matrices[0][centerCell.location.x][centerCell.location.y].state = "C"

	//Updating generations of matrices
	for m := 1; m <= numGens; m++ {
		fmt.Println("Updating " + strconv.Itoa(m) + "th generation...")
		matrices[m] = Update2DMatrix(matrices[m-1], Kcc, Knn, Knc)

	}

	return matrices
}

// Update2DMatrix takes in one matrix and coupling constants and returns a matrix of updated:
// 1) states (i.e., cancer cells can either stay proliferative, turn quiescent (and vice versa), or die) per lattice-gas/Boltzmann probability model.
// 2) velocities based on rules for necrotic of cancerous neighbors
func Update2DMatrix(currMatrix Matrix2D, Kcc, Knn, Knc float64) Matrix2D { //returns updated matrix (doesn't edit old one since we want to plot all!)

	x := GetNumRows2D(currMatrix)

	y := GetNumCols2D(currMatrix)

	//updating cell states based upon probabilities calculated using prior matrix. This will update each cell such that cells are currently cancerous.
	statesMatrix := UpdateMatrixStates2D(currMatrix, x, y, Kcc, Knn, Knc) //

	// updating cell velocities (transport step) based on rules for necrotic and cancerous cells in neighborhood.
	velocitiesMatrix := UpdateMatrixVelocities2D(statesMatrix, x, y) //, freshMatrix

	//now that states and velocities at each cell are reflected in freshMatrix, it's time to proliferate cells if applicable.
	readyMatrix := Initialize2DMatrix(x, y)

	readyMatrix = velocitiesMatrix

	//and push the cells according to the pushing rules
	pushedMatrix := PushAllCells2D(readyMatrix, x, y)

	return pushedMatrix
}

//UpdateMatrixStates2D updates the states of the cells using UpdateOneCellState2D subroutine
func UpdateMatrixStates2D(currMatrix Matrix2D, x, y int, Kcc, Knn, Knc float64) Matrix2D {

	statesMatrix := Initialize2DMatrix(x, y)

	for i := range currMatrix {

		for j := range currMatrix[i] {

			if InField2D(i, j, x, y) == true {
				statesMatrix[i][j] = currMatrix[i][j]

				//if cell is a living cancer cell, we update to C (will propagate/proliferate), or Q (quiescent; still alive, but will not progagate), or N (cell dies.)
				if currMatrix[i][j].state == "C" || currMatrix[i][j].state == "Q" {

					// updating cell states in new matrix based on current states (of prev matrix)
					statesMatrix[i][j] = UpdateOneCellState2D(currMatrix, i, j, Kcc, Knn, Knc)
				}
			}
		}
	}

	return statesMatrix
}

//UpdateOneCellState2D updates cell states based on previous probabilities and adds current probabilities to Cell2D struct
func UpdateOneCellState2D(currMatrix Matrix2D, i, j int, Kcc, Knn, Knc float64) Cell2D {

	numRows := GetNumRows2D(currMatrix)
	numCols := GetNumCols2D(currMatrix)

	currNhd := GetCurrentNeighborhood2D(currMatrix, i, j, numRows, numCols)

	C := GetNumCancerous2D(currNhd) //want to discount center cell.
	N := GetNumNecrotic2D(currNhd)

	Ep := EProliferation(Kcc, Knn, Knc, N, C)
	Eq := EQuiescence(Kcc, Knn, Knc, N, C)
	En := ENecrosis(Kcc, Knn, Knc, N, C)

	pN := ProbNecrosis(Ep, En, Eq)

	//Probability had to be edited to proliferate. Without the multiplication, the pP and pQ values were too low.
	pP := ProbProliferation(Ep, En, Eq) * 100000000.0
	pQ := ProbQuiescence(Ep, En, Eq) * 100000.0

	var newCell Cell2D

	newCell = currMatrix[i][j]

	currMatrix[i][j].pNecrosis = pN
	currMatrix[i][j].pQuiescent = pQ
	currMatrix[i][j].pProliferation = pP

	pAll := []float64{pN, pP, pQ}

	//getting max of probabilities for next state
	maxP := GetMaxP(pAll)

	//traceback step:
	//quiescent, necrotic, and cancerous are possible next states
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

//GetMaxP retrieves max probability from a slice of probabilities.
func GetMaxP(allP []float64) float64 {

	maxP := 0.0
	for i := range allP {

		if allP[i] > maxP {

			maxP = allP[i]
		}
	}
	return maxP
}

//UpdateMatrixVelocities2D updates the Matrix2D by utilizing subroutine that updates one cell
func UpdateMatrixVelocities2D(statesMatrix Matrix2D, x, y int) Matrix2D {

	velocitiesMatrix := Initialize2DMatrix(x, y)

	numRows := GetNumRows2D(statesMatrix)
	numCols := GetNumCols2D(statesMatrix)

	for i := range statesMatrix {

		for j := range statesMatrix[i] {

			if InField2D(i, j, numRows, numCols) == true {

				velocitiesMatrix[i][j] = UpdateOneCellVelocity2D(statesMatrix, i, j) // updating cell states
			}
		}
	}

	return velocitiesMatrix
}

//UpdateOneCellVelocity2D updates the velocity direction of a cell
func UpdateOneCellVelocity2D(statesMatrix Matrix2D, i, j int) Cell2D {
	//should indicate direction. magnitude assumed to be 1 for now.

	currCell := statesMatrix[i][j]

	numRows := GetNumRows2D(statesMatrix)
	numCols := GetNumCols2D(statesMatrix)

	if InField2D(i, j, numRows, numCols) == true {

		//if PROLIFERATIVE cancerous cell (C), then velocity vector should point at the direction of least (C+Q) cells..
		if currCell.state == "C" {

			minCNeighborCoord := GetMinCNeighborDirection2D(statesMatrix, i, j)

			//set velocity vector to point to direction of neighbor least-dense with cancer cells.
			currCell.velocityDirection = minCNeighborCoord

		} else if currCell.state == "N" {

			maxNNeighborCoord := GetMaxNNeighborDirection2D(statesMatrix, i, j)

			//set velocity vector to point to neighbor with most necrosis in its neighborhood.
			currCell.velocityDirection = maxNNeighborCoord

		} else {
			currCell.velocityDirection.x = currCell.location.x
			currCell.velocityDirection.x = currCell.location.y
		}
	}
	return currCell
}

//GetMaxNNeighborDirection2D retreives coordinates of cell neighbor to i,j with its neighborhood as most-dense in N.
func GetMaxNNeighborDirection2D(currMatrix Matrix2D, i, j int) OrderedPair {

	numRows := GetNumRows2D(currMatrix)
	numCols := GetNumCols2D(currMatrix)

	//ranging over all neighboring cells to see which neighbor is itself surrounded by fewest cancer cells.

	//Noah Chang------------------------------------------------------------------
	//Debugged here for handling the edges
	neighborhoods := make([]Neighborhood2D, 0)
	if InField2D(i+2, j, numRows, numCols) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood2D(currMatrix, i+2, j, numRows, numCols))
	}
	if InField2D(i-2, j, numRows, numCols) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood2D(currMatrix, i-2, j, numRows, numCols))
	}
	if InField2D(i, j+2, numRows, numCols) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood2D(currMatrix, i, j+2, numRows, numCols))
	}
	if InField2D(i, j-2, numRows, numCols) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood2D(currMatrix, i, j-2, numRows, numCols))
	}
	//----------------------------------------------------------------------------

	var maxNcoords OrderedPair

	maxCountN := 0.0 // zero the necrotic count since we want the max. Necrotic cells are chemotactic to others.

	for i := range neighborhoods { //ranging over surrounding nhds.

		currCountN := GetNumNecrotic2D(neighborhoods[i])

		if currCountN > maxCountN { //if a new MINIMUM is found, set coordinates (this is where we WANT a cancerous cell to go).

			maxNcoords.x = neighborhoods[i].neighbors[rand.Intn(len(neighborhoods))].location.x
			maxNcoords.y = neighborhoods[i].neighbors[rand.Intn(len(neighborhoods))].location.y

		} else if currCountN == maxCountN {

			maxCountN = currCountN // TIEBREAKING: we take at random the maximum N count.

			maxNcoords.x = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.x
			maxNcoords.y = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.y

		}
	}

	return maxNcoords

}

//GetMinCNeighborDirection2D returns coordinates of neighbor of cell with minimum cancer density by ranging over the neighborhoods OF the neighbors to the given cell at position i,j
func GetMinCNeighborDirection2D(currMatrix Matrix2D, i, j int) OrderedPair {

	numRows := GetNumRows2D(currMatrix)
	numCols := GetNumCols2D(currMatrix)

	//ranging over all neighboring cells to see which surrounding neighborhood contains fewest cancer cells.

	var minCcoords OrderedPair

	minCountC := float64(numRows * numCols) // make a large number out of cubed dimensions of matrix

	//Noah Chang------------------------------------------------------------------
	neighborhoods := make([]Neighborhood2D, 0)
	if InField2D(i+2, j, numRows, numCols) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood2D(currMatrix, i+2, j, numRows, numCols))
	}
	if InField2D(i-2, j, numRows, numCols) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood2D(currMatrix, i-2, j, numRows, numCols))
	}
	if InField2D(i, j+2, numRows, numCols) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood2D(currMatrix, i, j+2, numRows, numCols))
	}
	if InField2D(i, j-2, numRows, numCols) == true {
		neighborhoods = append(neighborhoods, GetCurrentNeighborhood2D(currMatrix, i, j-2, numRows, numCols))
	}
	//----------------------------------------------------------------------------

	for i := range neighborhoods { //ranging over surrounding nhds.

		currCountC := GetNumCancerous2D(neighborhoods[i])

		if currCountC < minCountC { //if a new MINIMUM is found, set coordinates (this is where we WANT a cancerous cell to go).

			minCcoords.x = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods))].location.x
			minCcoords.y = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods))].location.y

		} else if currCountC == minCountC {

			minCountC = currCountC // TIEBREAKING: we take at random the maximum N count.

			minCcoords.x = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.x
			minCcoords.y = neighborhoods[rand.Intn(len(neighborhoods))].neighbors[rand.Intn(len(neighborhoods[i].neighbors))].location.y

		}

	}

	return minCcoords

}

//PushAllCells2D pushes a new cell to relevent coordinate given by the velocity neighborhood
func PushAllCells2D(currMatrix Matrix2D, x, y int) Matrix2D {

	pushedMatrix := Initialize2DMatrix(x, y)

	for i := range currMatrix {

		for j := range currMatrix[i] {

			pushedMatrix[i][j] = currMatrix[i][j]

		}
	}

	//  if the cells have their state as proliferative ("C"), we push a new cell to the relevant coordinate given by the velocity neighborhood. (COPYING "C" state to new position)
	//	if a cell is quiescent ()"Q"), do nothing.
	// 	if a cell is necrotic, MOVE it to most necrotic direction (per velocity of that cell) (deleting previous position)

	//first, let's update the matrix (copied) and push necrotic cells (less motile so more realistic to do this versus synchronously with cancer.)

	for i := range currMatrix {

		for j := range currMatrix[i] {

			currCell := currMatrix[i][j]

			toX := currCell.velocityDirection.x
			toY := currCell.velocityDirection.y

			if currCell.state == "C" {
				//distributing cancer cells
				//keep old cancer cell at original location and replace cell state at target.
				pushedMatrix[i][j].state = "C"
				pushedMatrix[toX][toY].state = "C" //cancer cell proliferates, but original cancer cell persists. Quiescent cells have no change.
			}

			if currCell.state == "N" {
				//necrotic cells move toward necrotic cells. Cancer cells will move toward non-cancer cells (normal and necrotic)

				pushedMatrix[i][j].state = "wN"    // blank since idea is that necrotic cell moved away from original position.
				pushedMatrix[toX][toY].state = "N" // "move" cell to location of vector pointer
			}

		}
	}

	return pushedMatrix

}

//LatticeConfigEnergy2D sums energy over all neighborhoods
//This function may be used for improving the model
func LatticeConfigEnergy2D(currMatrix Matrix2D, Kcc, Knn, Knc float64) float64 {

	latticeConfigEnergy := 0.0 //sum of energy over all neighborhoods

	numRows := GetNumRows2D(currMatrix)
	numCols := GetNumCols2D(currMatrix)

	for i := range currMatrix {
		for j := range currMatrix[i] {

			currNeighborhood := GetCurrentNeighborhood2D(currMatrix, i, j, numRows, numCols)

			currNhdEnergy := NeighborhoodConfigEnergy2D(currNeighborhood, Kcc, Knn, Knc)

			latticeConfigEnergy += currNhdEnergy

		}
	}

	//for all neighborhoods in lattice, calculate configuration energy.

	return latticeConfigEnergy
}

//Initialize2DMatrix makes a board full of cells without seeding
func Initialize2DMatrix(x, y int) Matrix2D {

	matrix := make(Matrix2D, x)

	for i := range matrix {

		matrix[i] = make([]Cell2D, y)

	}

	for l := range matrix { //looping through all cells and replacing nil vals.

		for m := range matrix[l] {

			matrix[l][m].state = "h" //healthy normal cells (boundary cases).

			var cellLocation OrderedPair
			cellLocation.x = l
			cellLocation.y = m

			matrix[l][m].location = cellLocation //assigning coordinates to track cell locations.

		}
	}

	return matrix
}

//GetCurrentNeighborhood2D returns slice of pointers to cells and 3D moore rulestring. x,y are current positions of center cell.
func GetCurrentNeighborhood2D(currMatrix Matrix2D, x, y int, numRows, numCols int) Neighborhood2D {

	currCenterCell := currMatrix[x][y]

	var currNhd Neighborhood2D //establishing a slice of cells that represent VonNeumann Neighborhood
	// first element will be (pointer to) center of neighborhood of first matrixes

	currNhd.center = &currCenterCell

	// MOORE VERSION: much more complex memory-wise and may not be as realistic
	if InField2D(x, y, numRows, numCols) == true { //discount center cell and make sure we are in the field.

		neighborSlice := []*Cell2D{
			&currMatrix[x-1][y],
			&currMatrix[x+1][y],
			&currMatrix[x][y-1],
			&currMatrix[x][y+1],
		}

		for _, nbr := range neighborSlice {

			currNhd.neighbors = append(currNhd.neighbors, nbr)

		}
	}

	return currNhd
}

//InField2D returns true if the given coordinate is in the field
func InField2D(x, y int, numRows, numCols int) bool {

	// since we check neighborhoods of neighbors of a given cell, the border case is TWO cells in magnitude

	if (x-5 < 0 || x+5 > numRows) || (y-5 < 0 || y+5 > numCols) {
		return false //out of matrix field
	}
	return true //in the matrix field.
}

//AssertSquareMatrix ranges over the matrix and ensures that they all have same length
func AssertSquareMatrix(currMatrix Matrix2D) {
	if len(currMatrix) == 0 {
		fmt.Println("Game board has no rows.")
		os.Exit(2)
	}

	numCols := len(currMatrix[0])
	for i := 1; i < len(currMatrix); i++ {
		if len(currMatrix[i]) != numCols {
			fmt.Println("Board isn't square.")
			os.Exit(1)
		}
	}
}

//GetNumCancerous2D retreives the number of cancerous cells in the neighborhood and also the center
func GetNumCancerous2D(nhd Neighborhood2D) float64 {

	numC := 0.0 //initialize to zero

	for i := range nhd.neighbors {
		if nhd.neighbors[i].state == "C" || nhd.neighbors[i].state == "Q" {
			numC++
		}
	}
	//now, center cell imputed
	if nhd.center.state == "C" || nhd.center.state == "Q" { //"proliferative" || nhd.center.state == "quiescent" {
		numC++
	}

	return numC
}

//GetNumNecrotic2D , in a given neighborhood, gets the number of Necrotic Cells.
func GetNumNecrotic2D(nhd Neighborhood2D) float64 {

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

//GetNumRows2D gets the number of rows
func GetNumRows2D(matrix Matrix2D) int {

	rows := len(matrix)

	return rows
}

//GetNumCols2D gets the number of columns
func GetNumCols2D(matrix Matrix2D) int {

	cols := len(matrix[0])

	return cols
}

//GetCellState2D gets the state of the cell
func GetCellState2D(currCell Cell2D) string {
	state := currCell.state
	return state
}
