package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/rcssggb/ggb-lib/rcsscommon"
	"github.com/rcssggb/ggb-lib/trainerclient"
	"github.com/rcssggb/ggb-single/player"
	q "github.com/rcssggb/ggb-single/qlearning"
	"gorgonia.org/tensor"
)

func main() {
	epsilon := 0.1
	const alpha = float32(0.2)
	const epsilonDecay = 0.9999
	naiveGames := 0
	gameCounter := 0
	weightsFileA := "weightsA.rln"
	weightsFileB := "weightsB.rln"

	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	hostName := "rcssserver"

	var qLearningA, qLearningB *q.QLearning
	var err error

	_, err = os.Stat(weightsFileA)
	if os.IsNotExist(err) {
		log.Println("creating new agent")
		qLearningA, err = q.Init()
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("loading agent from %s\n", weightsFileA)
		qLearningA, err = q.Load(weightsFileA)
	}

	_, err = os.Stat(weightsFileB)
	if os.IsNotExist(err) {
		log.Println("creating new agent")
		qLearningB, err = q.Init()
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("loading agent from %s\n", weightsFileB)
		qLearningB, err = q.Load(weightsFileB)
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
		startX, startY := rcsscommon.RandomPosition()
		if startX > 0 {
			startX = -startX
		}
		t.MovePlayer("single-agent", 1, startX, startY, 0, 0, 0)
		t.Start()
		// lastGoalTime := -1
		currentTime := 0
		returnValue := float32(0)
		for {
			if p.Client.PlayMode() == rcsscommon.PlayModeTimeOver {
				p.Client.Log(p.Client.Bye())
				break
			}

			// Choose A from S using policy derived from Q (e.g., epsilon-greedy)
			qValuesA, err := qLearningA.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}
			qValuesB, err := qLearningB.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}

			var qValues *tensor.Dense
			qSum, err := qValuesA.Add(qValuesB)
			if err != nil {
				p.Client.Log(err)
			}
			qValues, err = qSum.DivScalar(float32(2.0), true)
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
			r := float32(0)
			ball := p.GetBall()
			if ball.NotSeenFor == 0 {
				ballDist := float32(ball.Distance)
				if ballDist < 0.7 {
					ballDist = 0.7
					epsilon *= epsilonDecay
				}
				r = 1.0 / ballDist
			}
			returnValue += r
			// if p.Client.PlayMode() == rcsscommon.PlayModeGoalL && currentTime > lastGoalTime {
			// 	lastGoalTime = currentTime
			// 	r = 1
			// 	p.Client.Log("goal!")
			// 	epsilon = epsilon * epsilonDecay
			// }

			// Update Q towards target
			nextActionValuesA, err := qLearningA.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}
			nextActionValuesB, err := qLearningB.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}
			if rand.Float32() < 0.5 {
				maxActionCoord, err := nextActionValuesA.Argmax(1)
				if err != nil {
					p.Client.Log(err)
				}
				maxActionCoordVal := maxActionCoord.Data().([]int)[0]
				nextMax := nextActionValuesB.Get(maxActionCoordVal)
				if err != nil {
					p.Client.Log(err)
				}
				nextMaxVal := nextMax.(float32)
				if math.IsNaN(float64(nextMaxVal)) {
					panic("training diverged")
				}
				td := r + nextMaxVal
				currentQ := qValuesA.Get(action)
				currentQVal := currentQ.(float32)
				qValuesA.Set(action, currentQVal+alpha*(td-currentQVal))
				err = qLearningA.Update(state, qValuesA)
				if err != nil {
					p.Client.Log(err)
				}
			} else {
				maxActionCoord, err := nextActionValuesB.Argmax(1)
				if err != nil {
					p.Client.Log(err)
				}
				maxActionCoordVal := maxActionCoord.Data().([]int)[0]
				nextMax := nextActionValuesA.Get(maxActionCoordVal)
				if err != nil {
					p.Client.Log(err)
				}
				nextMaxVal := nextMax.(float32)
				if math.IsNaN(float64(nextMaxVal)) {
					panic("training diverged")
				}
				td := r + nextMaxVal
				currentQ := qValuesB.Get(action)
				currentQVal := currentQ.(float32)
				qValuesB.Set(action, currentQVal+alpha*(td-currentQVal))
				err = qLearningB.Update(state, qValuesB)
				if err != nil {
					p.Client.Log(err)
				}
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
			err = qLearningA.Save(weightsFileA)
			if err != nil {
				log.Println(err)
			}
			err = qLearningB.Save(weightsFileB)
			if err != nil {
				log.Println(err)
			}
			log.Printf("weights saved after %d games\n", gameCounter)
		}
		time.Sleep(5 * time.Second)
	}
}
