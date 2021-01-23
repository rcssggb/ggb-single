package qlearning

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"

	"github.com/aunum/goro/pkg/v1/layer"
	m "github.com/aunum/goro/pkg/v1/model"
	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// QLearning contains sequential model data for training
type QLearning struct {
	actionValue *m.Sequential
}

// InitQLearning instantiates models for actor-critic learning of
// a discrete policy
func InitQLearning(stateSize, actionSize int) (*QLearning, error) {
	qModel, err := m.NewSequential("qLearning")
	if err != nil {
		return nil, err
	}

	inSize := stateSize + actionSize

	xShape := []int{1, inSize}
	yShape := []int{1, 1}

	qModel.AddLayers(
		layer.FC{Input: inSize, Output: 64},
		layer.FC{Input: 256, Output: 64},
		layer.FC{Input: 64, Output: 16},
		layer.FC{Input: 4, Output: 1, Activation: layer.Linear},
	)

	in := m.NewInput("state", xShape)
	out := m.NewInput("actionValue", yShape)

	fmt.Println("in shape: ", in.Shape())

	err = qModel.Compile(in, out,
		m.WithBatchSize(1),
	)
	if err != nil {
		return nil, err
	}

	return &QLearning{actionValue: qModel}, nil
}

// Slice2Tensor converts slice-formatted state into gorgonia tensor
func Slice2Tensor(state []float64) *tensor.Dense {
	return tensor.New(tensor.WithShape(1, len(state)), tensor.WithBacking(state))

}

// GreedyAction TODO: we need to append each action vector to state vector before predicting
func (q *QLearning) GreedyAction(state *tensor.Dense) (greedyA int, greedyValue float64, err error) {
	greedyA = 0
	greedyValue = -math.MaxFloat64
	for a := 0; a < 16; a++ {
		var val gorgonia.Value
		val, err = q.actionValue.Predict(state)
		if err != nil {
			return
		}

		valNum, ok := val.Data().(float64)
		fmt.Println(reflect.TypeOf(val.Data()), val.Data())
		if !ok {
			err = fmt.Errorf("val.Data() is not as you thought it was")
			return
		}

		if valNum > greedyValue {
			greedyValue = valNum
			greedyA = a
		}
	}
	return
}

// SampleDiscreteActionVector samples a random discrete action vector
func SampleDiscreteActionVector() (int, []float64) {
	a := rand.Intn(16)
	vec := make([]float64, 16)
	vec[a] = 1
	return a, vec
}

// DiscreteActionVector generates an action vector
func DiscreteActionVector(a int) []float64 {
	vec := make([]float64, 16)
	vec[a] = 1
	return vec
}
