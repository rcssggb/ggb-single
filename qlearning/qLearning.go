package qlearning

import (
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
		layer.FC{Input: in.Squeeze()[0], Output: 256, Init: gorgonia.Zeroes()},
		layer.FC{Input: 256, Output: 128, Init: gorgonia.Zeroes()},
		layer.FC{Input: 128, Output: 64, Init: gorgonia.Zeroes()},
		layer.FC{Input: 64, Output: 32, Init: gorgonia.Zeroes()},
		layer.FC{Input: 32, Output: out.Squeeze()[0], Activation: layer.Linear, Init: gorgonia.Zeroes()},
	)

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

// ActionValues returns action values for all possible actions
func (q *QLearning) ActionValues(state *tensor.Dense) (actionValues *tensor.Dense, err error) {
	var val gorgonia.Value
	val, err = q.actionValue.Predict(state)
	if err != nil {
		return
	}
	actionValues = val.(*tensor.Dense)

	return
}

// Update updates learnables from state towards target
func (q *QLearning) Update(state, target *tensor.Dense) error {
	return q.actionValue.Fit(state, target)
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
