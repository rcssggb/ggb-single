package linearq

import (
	"gonum.org/v1/gonum/mat"
)

// LinearQ contains linear approximator data for training
type LinearQ struct {
	state2ActionValue mat.Matrix
}

// Init initializes LinearQ object
func Init(stateSize, actionSize int) *LinearQ {
	weights := make([]float64, actionSize*(stateSize+1))
	wMat := mat.NewDense(actionSize, stateSize+1, weights)
	return &LinearQ{state2ActionValue: wMat}
}

// Slice2Vector converts slice-formatted state into gorgonia tensor
func Slice2Vector(state []float64) mat.Vector {
	sVec := mat.NewVecDense(len(state)+1, append(state, 1))
	return sVec
}
