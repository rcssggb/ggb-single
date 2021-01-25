package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/rcssggb/ggb-lib/rcsscommon"
	"github.com/rcssggb/ggb-lib/trainerclient"
	"github.com/rcssggb/ggb-single/player"
	q "github.com/rcssggb/ggb-single/qlearning"
)

func main() {
	epsilon := 0.9
	const epsilonDecay = 0.999

	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	hostName := "rcssserver"

	qLearning, err := q.InitQLearning(71, 16)
	if err != nil {
		panic(err)
	}

	for {
		p, err := player.NewPlayer("single-agent", hostName)
		if err != nil {
			log.Println(err)
			continue
		}

		t, err := trainerclient.NewTrainerClient(hostName)
		if err != nil {
			log.Println(err)
			continue
		}

		t.EarOn()
		t.EyeOn()

		time.Sleep(2 * time.Second)

		// Initialize S
		state := q.Slice2Tensor(p.State())
		t.Start()
		lastGoalTime := -1
		for {
			if p.Client.PlayMode() == rcsscommon.PlayModeTimeOver {
				p.Client.Log(p.Client.Bye())
				break
			}

			currentTime := p.Client.Time()

			// Choose A from S using policy derived from Q (e.g., epsilon-greedy)
			qValues, err := qLearning.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}
			var action int
			takeRandomAction := rand.Float64() < epsilon
			if takeRandomAction {
				action = rand.Intn(16)
			} else {
				maxActionTensor, err := qValues.Argmax(1)
				if err != nil {
					p.Client.Log(err)
				}
				action = maxActionTensor.Data().(int)
			}

			// Take action A
			p.DiscreteAction(action)

			err = p.Client.Error()
			for err != nil {
				p.Client.Log(err)
				err = p.Client.Error()
			}

			if currentTime != 0 {
				// time.Sleep(100 * time.Millisecond)
			}
			t.DoneSynch()
			p.WaitCycle()

			// Observe R, S'
			nextState := q.Slice2Tensor(p.State())
			var r float32
			if p.Client.PlayMode() == rcsscommon.PlayModeGoalL && currentTime > lastGoalTime {
				lastGoalTime = currentTime
				r = 1
				p.Client.Log("goal!")
				epsilon = epsilon * epsilonDecay
			}

			// Update Q towards target
			nextActionValues, err := qLearning.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}
			nextMax, err := nextActionValues.Max(1)
			if err != nil {
				p.Client.Log(err)
			}
			td := r + nextMax.Data().(float32)
			qValues.Set(action, td)
			err = qLearning.Update(state, qValues)
			if err != nil {
				p.Client.Log(err)
			}

			// S <- S'
			state = nextState
		}
	}
}
