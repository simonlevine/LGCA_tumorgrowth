package main

// The following code was written by Simon Levine-Gottreich

//NeighborhoodConfigEnergy2D calculates neighborhood config energy
func NeighborhoodConfigEnergy2D(currNeighborhood Neighborhood2D, Kcc, Knn, Knc float64) float64 {

	C := GetNumCancerous2D(currNeighborhood) //includes center cell.
	N := GetNumNecrotic2D(currNeighborhood)

	// apoptosized cells and proliferative cells are not needed here...

	//by literature formula...

	Econfig := -1 * (.50*(C*(C-1)*Kcc+N*(N-1)*Knn) + C*N*Knc)

	return Econfig
}

// The following I will include for the sake of modularity: The delta functions allow for a more complex computational method involving physical constants.
// But, in the current implementation, we instead use proportional formulae (i.e, presence or absnce) for simplicity and to remain within 64-bit floating point precision.

// func DeltaEQuiescence2D(currNeighborhood Neighborhood2D, Kcc, Knn, Knc float64) float64 { //should be zero, since no change.
// 	//want all neighbors, not center cell.
//
// 	currNhdEnergy := NeighborhoodConfigEnergy2D(currNeighborhood, Kcc, Knn, Knc)
//
// 	//no change to number cancer or necrotic cells.
//
// 	energyIfQuiescent := currNhdEnergy
//
// 	deltaEQ := energyIfQuiescent - currNhdEnergy
//
// 	return deltaEQ
// }
//
// func DeltaEApoptosis2D(currNeighborhood Neighborhood2D, Kcc, Knn, Knc float64) float64 {
//
// 	// if and only if (C ≥1)
//
// 	currNhdEnergy := NeighborhoodConfigEnergy2D(currNeighborhood, Kcc, Knn, Knc)
//
// 	// apoptosized cells and proliferative cells are not needed here...
//
// 	//by literature formula... Remove a cancer cell.
//
// 	C := GetNumCancerous2D(currNeighborhood) - 1.0 //1 dead cancer cell.
// 	N := GetNumNecrotic2D(currNeighborhood)
//
// 	energyIfApoptotic := -1 * (.50*(C*(C-1)*Kcc+N*(N-1)*Knn) + C*N*Knc)
//
// 	deltaEA := energyIfApoptotic - currNhdEnergy
//
// 	return deltaEA
// }
//
// func DeltaEProliferation2D(currNeighborhood Neighborhood2D, Kcc, Knn, Knc float64) float64 {
//
// 	currNhdEnergy := NeighborhoodConfigEnergy2D(currNeighborhood, Kcc, Knn, Knc)
//
// 	C := GetNumCancerous2D(currNeighborhood) + 1.0 //1 MORE cancer cell.
// 	N := GetNumNecrotic2D(currNeighborhood)
//
// 	energyIfProliferative := -1 * (.50*(C*(C-1)*Kcc+N*(N-1)*Knn) + C*N*Knc)
//
// 	deltaEP := energyIfProliferative - currNhdEnergy
//
// 	return deltaEP
// }
//
// func DeltaENecrosis2D(currNeighborhood Neighborhood2D, Kcc, Knn, Knc float64) float64 {
//
// 	currNhdEnergy := NeighborhoodConfigEnergy2D(currNeighborhood, Kcc, Knn, Knc)
//
// 	C := GetNumCancerous2D(currNeighborhood) - 1.0 //1 LESS cancer cell.
// 	N := GetNumNecrotic2D(currNeighborhood) + 1.0  //1 MORE dead cell.
//
// 	energyIfNecrotic := -1 * (.50*(C*(C-1)*Kcc+N*(N-1)*Knn) + C*N*Knc)
//
// 	deltaEN := energyIfNecrotic - currNhdEnergy
//
// 	return deltaEN
//
// }
//

// func ProbApoptosis(deltaEP, deltaEA, deltaEN, deltaEQ float64) float64 {
//
// 	bottom := math.Exp(deltaEQ) + math.Exp(deltaEP) + math.Exp(deltaEA) + math.Exp(deltaEN)
// 	//setting denominator
//
// 	top := math.Exp(deltaEA) // setting numerator
//
// 	pApoptosis := top / bottom
//
// 	return pApoptosis
// }
//
// remember cases when boltzmann constants are 0.
// //must take in number cancerous, number necrotic (dead)
// func ProbNecrosis(deltaEP, deltaEA, deltaEN, deltaEQ float64) float64 {
//
// 	bottom := math.Exp(deltaEQ) + math.Exp(deltaEP) + math.Exp(deltaEN) //+ math.Exp(deltaEA)
// 	//setting denominator
//
// 	top := math.Exp(deltaEN) // setting numerator
//
// 	pNecrosis := top / bottom
//
// 	return pNecrosis
//
// }
//
// func ProbProliferation(deltaEP, deltaEA, deltaEN, deltaEQ float64) float64 { //N, C float64, k, T float64 (add in later?)
//
// 	// ***SIMPLIFYING MODEL ***
//
// 	bottom := math.Exp(deltaEQ) + math.Exp(deltaEP) + math.Exp(deltaEN) //+ math.Exp(deltaEA)
// 	//setting denominator
//
// 	top := math.Exp(deltaEP) // setting numerator
//
// 	pProliferation := (top / bottom)
//
// 	return pProliferation
// }
//
// func ProbQuiescence(deltaEP, deltaEA, deltaEN, deltaEQ float64) float64 {
//
// 	bottom := math.Exp(deltaEQ) + math.Exp(deltaEP) + math.Exp(deltaEN) //+ math.Exp(deltaEA)
// 	//setting denominator
//
// 	top := math.Exp(deltaEQ) // setting numerator
//
// 	pQuiescence := top / bottom
//
// 	return pQuiescence
// }
