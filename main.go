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

			// TODO: Choose A from S using policy derived from Q (e.g., epsilon-greedy)
			action, _, err := qLearning.GreedyAction(state)
			if err != nil {
				p.Client.Log(err)
			}

			// Take action A
			p.DiscreteAction(action)

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
			var r float64
			if p.Client.PlayMode() == rcsscommon.PlayModeGoalL {
				r = 1
				if r == 1 {
					p.Client.Log("goal!")
				}
			}

			// TODO: Update Q towards target

			// S <- S'
			state = nextState
		}
	}
}
