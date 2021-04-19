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
	actionValue  *m.Sequential
	batchSize    int
	batchStates  []*tensor.Dense
	batchTargets []*tensor.Dense
	NStates      int
	NBehaviors   int
}

// Init instantiates models for q-learning of
// a discrete policy
func Init() (*QLearning, error) {
	q := QLearning{}
	q.NStates = 5
	q.NBehaviors = 8
	q.batchSize = 1
	q.batchStates = []*tensor.Dense{}
	q.batchTargets = []*tensor.Dense{}

	qModel, err := m.NewSequential("qLearning")
	if err != nil {
		return nil, err
	}

	xShape := []int{1, q.NStates}
	yShape := []int{1, q.NBehaviors}
	in := m.NewInput("state", xShape)
	out := m.NewInput("actionValue", yShape)

	qModel.AddLayers(
		layer.FC{Input: in.Squeeze()[0], Output: 2048, Init: gorgonia.GlorotN(0.001), BiasInit: gorgonia.GlorotN(0.001)},
		layer.FC{Input: 2048, Output: 1024, Init: gorgonia.GlorotN(0.001), BiasInit: gorgonia.GlorotN(0.001)},
		layer.FC{Input: 1024, Output: 512, Init: gorgonia.GlorotN(0.001), BiasInit: gorgonia.GlorotN(0.001)},
		layer.FC{Input: 512, Output: 256, Init: gorgonia.GlorotN(0.001), BiasInit: gorgonia.GlorotN(0.001)},
		layer.FC{Input: 256, Output: 128, Init: gorgonia.GlorotN(0.001), BiasInit: gorgonia.GlorotN(0.001)},
		layer.FC{Input: 128, Output: 64, Init: gorgonia.GlorotN(0.001), BiasInit: gorgonia.GlorotN(0.001)},
		layer.FC{Input: 64, Output: out.Squeeze()[0], Activation: layer.Linear, Init: gorgonia.GlorotN(0.001), BiasInit: gorgonia.GlorotN(0.001)},
	)
	err = qModel.Compile(in, out,
		m.WithBatchSize(q.batchSize),
		m.WithOptimizer(
			gorgonia.NewVanillaSolver(
				gorgonia.WithLearnRate(1),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	q.actionValue = qModel

	return &q, nil
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

// UpdateWithBatch updates learnables from state towards target
func (q *QLearning) UpdateWithBatch(state, target *tensor.Dense) error {
	q.batchStates = append(q.batchStates, state)
	q.batchTargets = append(q.batchTargets, target)
	if len(q.batchStates) >= q.batchSize {
		states, err := q.batchStates[0].Concat(0, q.batchStates[1:]...)
		if err != nil {
			return err
		}
		targets, err := q.batchTargets[0].Concat(0, q.batchTargets[1:]...)
		if err != nil {
			return err
		}
		q.batchStates = []*tensor.Dense{}
		q.batchTargets = []*tensor.Dense{}
		err = q.actionValue.FitBatch(states, targets)
		if err != nil {
			return err
		}
	}

	return nil
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
