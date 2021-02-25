package tabularq

import (
	"encoding/gob"
	"math"
	"math/rand"
	"os"
)

type actionValues []float64

func (a actionValues) Max() (int, float64) {
	max := -math.MaxFloat64
	action := 0
	for i := range a {
		if a[i] > max {
			action = i
			max = a[i]
		}
	}
	return action, max
}

type QLearning struct {
	A          []actionValues
	B          []actionValues
	NumStates  int
	NumActions int

	Alpha float64
	Gamma float64
}

// Init instantiates the QTable
func Init(nStates, nActions int) *QLearning {
	qTableA := make([]actionValues, nStates)
	for i := range qTableA {
		qTableA[i] = make(actionValues, nActions)
		for j := range qTableA[i] {
			qTableA[i][j] = rand.NormFloat64() * 1000
		}
	}

	qTableB := make([]actionValues, nStates)
	for i := range qTableB {
		qTableB[i] = make(actionValues, nActions)
		for j := range qTableB[i] {
			qTableB[i][j] = rand.NormFloat64() * 1000
		}
	}
	return &QLearning{
		A:          qTableA,
		B:          qTableB,
		NumStates:  nStates,
		NumActions: nActions,
		Alpha:      0.1, // default value
		Gamma:      1,
	}
}

// GreedyAction returns the greedy action according to q tables
func (q *QLearning) GreedyAction(state int) int {
	a := q.A[state]
	b := q.B[state]

	max := -math.MaxFloat64
	action := 0
	for i := 0; i < q.NumActions; i++ {
		sum := a[i] + b[i]
		if sum > max {
			action = i
			max = sum
		}
	}
	return action
}

// Update updates the q tables with better q approximations
func (q *QLearning) Update(state, action int, reward float64, nextState int) {
	var mainTable []actionValues
	var altTable []actionValues
	if rand.Intn(2) == 0 {
		mainTable = q.A
		altTable = q.B
	} else {
		mainTable = q.B
		altTable = q.A
	}

	nextGreedyAction, _ := altTable[state].Max()
	mainTable[state][action] = mainTable[state][action] + q.Alpha*(reward+q.Gamma*mainTable[nextState][nextGreedyAction]-mainTable[state][action])
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
