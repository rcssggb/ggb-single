package qlearning

import (
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/aunum/goro/pkg/v1/layer"
	m "github.com/aunum/goro/pkg/v1/model"
	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// QLearning contains sequential model data for training
type QLearning struct {
	actionValue *m.Sequential
}

// Init instantiates models for q-learning of
// a discrete policy
func Init() (*QLearning, error) {
	qModel, err := m.NewSequential("qLearning")
	if err != nil {
		return nil, err
	}

	xShape := []int{1, 71}
	yShape := []int{1, 16}
	in := m.NewInput("state", xShape)
	out := m.NewInput("actionValue", yShape)

	qModel.AddLayers(
		layer.FC{Input: in.Squeeze()[0], Output: 256, Init: gorgonia.Zeroes(), BiasInit: gorgonia.Zeroes()},
		layer.FC{Input: 256, Output: 128, Init: gorgonia.Zeroes(), BiasInit: gorgonia.Zeroes()},
		layer.FC{Input: 128, Output: 64, Init: gorgonia.Zeroes(), BiasInit: gorgonia.Zeroes()},
		layer.FC{Input: 64, Output: 32, Init: gorgonia.Zeroes(), BiasInit: gorgonia.Zeroes()},
		layer.FC{Input: 32, Output: out.Squeeze()[0], Activation: layer.Linear, Init: gorgonia.Zeroes(), BiasInit: gorgonia.Zeroes()},
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

// Save saves model
func (q *QLearning) Save(filename string) error {
	nodes := q.actionValue.Learnables()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	enc := gob.NewEncoder(f)
	for _, node := range nodes {
		err := enc.Encode(node.Value())
		if err != nil {
			return err
		}
	}
	return nil
}

// Load loads model
func Load(filename string) (*QLearning, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)

	i := 0
	q, err := Init()
	if err != nil {
		panic(err)
	}

	learnables := q.actionValue.Learnables()
	for {
		var t *tensor.Dense
		err = dec.Decode(&t)

		// Reach end of file
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if t == nil {
			fmt.Println("nil t")
			continue
		}

		learnable := learnables[i]
		err = gorgonia.Let(learnable, t)

		i++
	}

	q.actionValue.SetLearnables(learnables)

	return q, nil
}

// Learnables returns learnables
func (q *QLearning) Learnables() gorgonia.Nodes {
	return q.actionValue.Learnables()
}
