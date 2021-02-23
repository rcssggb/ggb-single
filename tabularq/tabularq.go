package tabularq

import (
	"math"
	"math/rand"
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
	a        []actionValues
	b        []actionValues
	nStates  int
	nActions int

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
		a:        qTableA,
		b:        qTableB,
		nStates:  nStates,
		nActions: nActions,
		Alpha:    0.1, // default value
		Gamma:    1,
	}
}

// GreedyAction returns the greedy action according to q tables
func (q *QLearning) GreedyAction(state int) int {
	a := q.a[state]
	b := q.b[state]

	max := -math.MaxFloat64
	action := 0
	for i := 0; i < q.nActions; i++ {
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
		mainTable = q.a
		altTable = q.b
	} else {
		mainTable = q.b
		altTable = q.a
	}

	nextGreedyAction, _ := altTable[state].Max()
	mainTable[state][action] = mainTable[state][action] + q.Alpha*(reward+q.Gamma*mainTable[nextState][nextGreedyAction]-mainTable[state][action])
}
