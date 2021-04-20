package linearq

import (
	"encoding/gob"
	"math"
	"math/rand"
	"os"

	"gonum.org/v1/gonum/mat"
)

type QLearning struct {
	A          *mat.Dense
	B          *mat.Dense
	NumStates  int
	NumActions int

	Alpha float64
	Gamma float64
}

// Init instantiates the weight matrices
func Init(nStates, nActions int) *QLearning {
	weights := make([]float64, nActions*nStates)
	for i := range weights {
		weights[i] = rand.NormFloat64()
	}
	qWeightsA := mat.NewDense(nActions, nStates, weights)

	for i := range weights {
		weights[i] = rand.NormFloat64()
	}
	qWeightsB := mat.NewDense(nActions, nStates, weights)
	return &QLearning{
		A:          qWeightsA,
		B:          qWeightsB,
		NumStates:  nStates,
		NumActions: nActions,
		Alpha:      0.1, // default value
		Gamma:      1,
	}
}

// GreedyAction returns the greedy action according to q tables
func (q *QLearning) GreedyAction(state mat.Vector) int {
	a := mat.NewDense(q.NumActions, 1, nil)
	a.Product(q.A, state)

	b := mat.NewDense(q.NumActions, 1, nil)
	b.Product(q.A, state)

	sum := mat.NewDense(q.NumActions, 1, nil)
	sum.Add(a, b)

	action := 0
	max := -math.MaxFloat64
	for i := 0; i < q.NumActions; i++ {
		v := sum.At(i, 0)
		if v > max {
			action = i
			max = v
		}
	}
	return action
}

// Update updates the q tables with better q approximations
func (q *QLearning) Update(state, action int, reward float64, nextState int) {
	var main *mat.Dense
	var alt *mat.Dense
	if rand.Intn(2) == 0 {
		main = q.A
		alt = q.B
	} else {
		main = q.B
		alt = q.A
	}

	// TODO write update rule
	// nextGreedyAction, _ := altTable[state].Max()
	// mainTable[state][action] = mainTable[state][action] + q.Alpha*(reward+q.Gamma*mainTable[nextState][nextGreedyAction]-mainTable[state][action])
}

// Save saves the Q table encoded as gob
func (q *QLearning) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(q)
	if err != nil {
		return err
	}
	return nil
}

// Load loads the Q table from file
func Load(filename string) (*QLearning, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	q := QLearning{}
	err = dec.Decode(&q)
	if err != nil {
		return nil, err
	}

	return &q, nil
}
