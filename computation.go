package main

import (
	"math"
)

// The following code was written by Simon Levine-Gottreich

//ProbNecrosis takes calculates the probability of necrosis
func ProbNecrosis(EP, EN, EQ float64) float64 {

	bottom := math.Exp(-1*EQ) + math.Exp(-1*EP) + math.Exp(-1*EN) //+ math.Exp(EA)
	//setting denominator

	top := math.Exp(-1 * EN) // setting numerator

	pNecrosis := top / bottom

	return pNecrosis

}

//ENecrosis computes the energy if necrotic
func ENecrosis(Kcc, Knn, Knc float64, N, C float64) float64 {

	energyIfNecrotic := -1 * (.50*((C-2)*(C-1)*Kcc+N*(N+1)*Knn) + (C-1)*(N+1)*Knc)

	EN := energyIfNecrotic ////changing from  config energy per Springer book (2014)

	return EN

}

//ProbProliferation calculates the probability of proliferation
func ProbProliferation(EP, EN, EQ float64) float64 {

	bottom := math.Exp(-1*EQ) + math.Exp(-1*EP) + math.Exp(-1*EN)
	//setting denominator

	top := math.Exp(EP) // setting numerator

	pProliferation := (top / bottom)

	return pProliferation
}

//EProliferation computes the energy if proliferative
func EProliferation(Kcc, Knn, Knc float64, N, C float64) float64 {

	energyIfProliferative := -1 * (.50*((C+1)*(C)*Kcc+N*(N-1)*Knn) + (C+1)*N*Knc)

	EP := energyIfProliferative

	return EP
}

//ProbQuiescence calculates the probability of quiescence.
func ProbQuiescence(EP, EN, EQ float64) float64 {

	bottom := math.Exp(-1*EQ) + math.Exp(-1*EP) + math.Exp(-1*EN) //setting denominator

	top := math.Exp(EQ) // setting numerator

	pQuiescence := top / bottom

	return pQuiescence
}

//EQuiescence computes the energy if quiescent
func EQuiescence(Kcc, Knn, Knc float64, N, C float64) float64 {

	energyIfQuiescent := -1 * (.50*(C*(C-1)*Kcc+N*(N-1)*Knn) + C*N*Knc)

	//currNhdEnergy + 0

	EQ := energyIfQuiescent

	return EQ
}

//NeighborhoodConfigEnergy calculate neighborhood config energy
func NeighborhoodConfigEnergy(currNeighborhood Neighborhood, Kcc, Knn, Knc float64, N, C float64) float64 {

	//by literature formula...

	Econfig := -1 * (.50*(C*(C-1)*Kcc+N*(N-1)*Knn) + C*N*Knc)

	return Econfig
}

//The following could be used for apoptotic modeling. Function works, but model now is simpler version.

// func ProbApoptosis(EP, EA, EN, EQ float64) float64 {
//
// 	bottom := math.Exp(EQ) + math.Exp(EP) + math.Exp(EA) + math.Exp(EN)
// 	//setting denominator
//
// 	top := math.Exp(EA) // setting numerator
//
// 	pApoptosis := top / bottom
//
// 	return pApoptosis
// }

//NOT currently implemented...
// func EApoptosis(Kcc, Knn, Knc float64, N, C float64) float64 {
//
// 	// if and only if (C ≥1)
//
// 	//currNhdEnergy := NeighborhoodConfigEnergy(currNeighborhood, Kcc, Knn, Knc)
//
// 	// apoptosized cells and proliferative cells are not needed here...
//
// 	//by literature formula... Remove a cancer cell.
//
// 	C = C - 1.0 //1 dead cancer cell.
//
// 	energyIfApoptotic := -1 * (.50*(C*(C-1)*Kcc+N*(N-1)*Knn) + C*N*Knc)
//
// 	EA := energyIfApoptotic
//
// 	return EA
// }
