package main

import (
	"log"
	"time"

	"github.com/rcssggb/ggb-lib/rcsscommon"
	"github.com/rcssggb/ggb-lib/trainerclient"
	"github.com/rcssggb/ggb-single/player"
	q "github.com/rcssggb/ggb-single/qlearning"
)

func main() {
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
		for {
			if p.Client.PlayMode() == rcsscommon.PlayModeTimeOver {
				p.Client.Log(p.Client.Bye())
				break
			}

			currentTime := p.Client.Time()

			// Choose A from S using policy derived from Q (e.g., epsilon-greedy)
			// TODO: implement epsilon-greedy behavior instead of pure greedy
			actionValues, err := qLearning.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}

			// Take action A
			maxActionTensor, err := actionValues.Argmax(1)
			if err != nil {
				p.Client.Log(err)
			}
			maxAction := maxActionTensor.Data().(int)
			p.DiscreteAction(maxAction)

			err = p.Client.Error()
			for err != nil {
				p.Client.Log(err)
				err = p.Client.Error()
			}

			if currentTime != 0 {
				time.Sleep(10 * time.Millisecond)
			}
			t.DoneSynch()
			p.WaitCycle()

			// Observe R, S'
			nextState := q.Slice2Tensor(p.State())
			var r float32
			if p.Client.PlayMode() == rcsscommon.PlayModeGoalL {
				r = 1
				if r == 1 {
					p.Client.Log("goal!")
				}
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
			actionValues.Set(maxAction, td)
			err = qLearning.Update(state, actionValues)
			if err != nil {
				p.Client.Log(err)
			}

			// S <- S'
			state = nextState
		}
	}
}
