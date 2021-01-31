package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/rcssggb/ggb-lib/rcsscommon"
	"github.com/rcssggb/ggb-lib/trainerclient"
	"github.com/rcssggb/ggb-single/player"
	q "github.com/rcssggb/ggb-single/qlearning"
)

func main() {
	epsilon := 0.0
	const epsilonDecay = 0.999
	naiveGames := 0
	gameCounter := 0
	weightsFile := "weights.rln"

	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	hostName := "rcssserver"

	var qLearning *q.QLearning
	var err error

	_, err = os.Stat(weightsFile)
	if os.IsNotExist(err) {
		log.Println("creating new agent")
		qLearning, err = q.Init()
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("loading agent from %s\n", weightsFile)
		qLearning, err = q.Load(weightsFile)
	}

	trainingStart := time.Now()
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
		currentTime := 0
		returnValue := float32(0)
		for {
			if p.Client.PlayMode() == rcsscommon.PlayModeTimeOver {
				p.Client.Log(p.Client.Bye())
				break
			}

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
				if naiveGames > 0 {
					action = p.NaivePolicy()
				} else {
					maxActionTensor, err := qValues.Argmax(1)
					if err != nil {
						p.Client.Log(err)
					}
					action = maxActionTensor.Data().([]int)[0]
				}
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
			currentTime = p.Client.Time()
			var r float32
			if p.Client.PlayMode() == rcsscommon.PlayModeGoalL && currentTime > lastGoalTime {
				lastGoalTime = currentTime
				r = 1
				p.Client.Log("goal!")
				epsilon = epsilon * epsilonDecay
			}
			returnValue += r

			// Update Q towards target
			nextActionValues, err := qLearning.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}
			nextMax, err := nextActionValues.Max(1)
			if err != nil {
				p.Client.Log(err)
			}
			nextMaxVal := nextMax.Data().([]float32)[0]
			if math.IsNaN(float64(nextMaxVal)) {
				panic("training diverged")
			}
			td := r + nextMaxVal
			qValues.Set(action, td)
			err = qLearning.Update(state, qValues)
			if err != nil {
				p.Client.Log(err)
			}

			// S <- S'
			state = nextState
		}
		if naiveGames > 0 {
			naiveGames--
		}
		gameCounter++
		timeSinceStart := time.Now().Sub(trainingStart)
		log.Printf("finished game number %d with return %f after training for %s with an average of %.1f seconds per game\n", gameCounter, returnValue, timeSinceStart, timeSinceStart.Seconds()/float64(gameCounter))
		if gameCounter%10 == 0 {
			err = qLearning.Save(weightsFile)
			if err != nil {
				fmt.Printf("weights saved after %d games\n", gameCounter)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
