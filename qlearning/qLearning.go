package qlearning

import (
	"fmt"
	"math"
	"math/rand"

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

	xShape := []int{1, stateSize}
	yShape := []int{1, actionSize}
	in := m.NewInput("state", xShape)
	out := m.NewInput("actionValue", yShape)

	qModel.AddLayers(
		layer.FC{Input: in.Squeeze()[0], Output: 256},
		layer.FC{Input: 256, Output: 128},
		layer.FC{Input: 128, Output: 64},
		layer.FC{Input: 64, Output: 32},
		layer.FC{Input: 32, Output: out.Squeeze()[0], Activation: layer.Linear},
	)

	fmt.Println("in shape: ", in.Shape())
	fmt.Println("out shape: ", out.Shape())

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
	f32state := make([]float32, len(state))
	for i, v := range state {
		f32state[i] = float32(v)
	}
	return tensor.New(tensor.WithShape(1, len(f32state)), tensor.WithBacking(f32state))
}

// GreedyAction chooses greedy action and returns which action it is and its value
func (q *QLearning) GreedyAction(state *tensor.Dense) (greedyA int, greedyValue float64, err error) {
	var val gorgonia.Value
	val, err = q.actionValue.Predict(state)
	if err != nil {
		return
	}
	valTensor, ok := val.Data().([][]float64)
	if !ok {
		err = fmt.Errorf("val.Data() is not as you thought it was")
		return
	}

	greedyA = 0
	greedyValue = -math.MaxFloat64
	for a := 0; a < 16; a++ {
		valNum := valTensor[0][a]
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
