package main

import (
	"encoding/gob"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/rcssggb/ggb-lib/rcsscommon"
	"github.com/rcssggb/ggb-lib/trainerclient"
	"github.com/rcssggb/ggb-single/player"
	q "github.com/rcssggb/ggb-single/qlearning"
)

func main() {
	epsilon := 0.7
	const K = float32(1)
	const gamma = float32(0.99)
	const epsilonDecay = 0.997
	naiveGames := 0
	gameCounter := 0
	weightsFile := "weights.rln"
	returnsFile := "./data/returns.rln"

	logName := time.Now().String() + ".log"
	file, err := os.OpenFile(path.Join("logs", logName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	log.Printf("starting training with // epsilon = %f // alpha = K / (K + 1) = %f // gamma = %f // epsilonDecay = %f // naiveGames = %d", epsilon, K/K, gamma, epsilonDecay, naiveGames)

	hostName := "rcssserver"

	var qLearning *q.QLearning

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

	returnValues := []float32{}
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
		p.Client.SynchSee()
		p.Client.ChangeView(rcsscommon.ViewWidthNarrow, rcsscommon.ViewQualityHigh)

		time.Sleep(2 * time.Second)

		// Initialize S
		state := q.Slice2Tensor(p.State())
		startX, startY := rcsscommon.RandomPosition()
		if startX > 0 {
			startX = -startX
		}
		t.MovePlayer("single-agent", 1, startX, startY, 0, 0, 0)
		t.Start()
		lastGoalTime := -1
		currentTime := 0
		returnValue := float32(0)
		for {
			// alpha is calculated so that
			// sum(alpha(currentTime)) = \inf; and
			// sum((alpha(currentTime)^2) < \inf
			// as noticed by Watkins and Dayan, 1992 https://link.springer.com/content/pdf/10.1007/BF00992698.pdf
			alpha := K / (K + float32(currentTime))
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
				action = rand.Intn(5)
			} else {
				if naiveGames > 0 {
					action = p.NaiveBehaviorPolicy()
				} else {
					maxActionTensor, err := qValues.Argmax(1)
					if err != nil {
						p.Client.Log(err)
					}
					action = maxActionTensor.Data().([]int)[0]
				}
			}
			// fmt.Println(action)
			p.ExecuteBehavior(action)

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
			lastTime := currentTime
			currentTime = p.Client.Time()

			// Wait until simulation cycle changes (1 action per simulation cycle)
			for currentTime == lastTime {
				t.DoneSynch()
				p.WaitCycle()
				lastTime = currentTime
				currentTime = p.Client.Time()
			}

			nextState := q.Slice2Tensor(p.State())
			currentTime = p.Client.Time()
			r := float32(0)

			if p.Client.PlayMode() == rcsscommon.PlayModeGoalL && currentTime > lastGoalTime {
				lastGoalTime = currentTime
				r = 1
				p.Client.Log("goal!")
			} else if p.Client.PlayMode() == rcsscommon.PlayModeGoalR && currentTime > lastGoalTime {
				lastGoalTime = currentTime
				r = -1
				p.Client.Log("goal against, bad!")
			}

			returnValue += r

			// Update Q towards target
			nextActionValues, err := qLearning.ActionValues(state)
			if err != nil {
				p.Client.Log(err)
			}
			maxActionCoord, err := nextActionValues.Argmax(1)
			if err != nil {
				p.Client.Log(err)
			}
			maxActionCoordVal := maxActionCoord.Data().([]int)[0]
			nextMax := nextActionValues.Get(maxActionCoordVal)
			if err != nil {
				p.Client.Log(err)
			}
			nextMaxVal := nextMax.(float32)
			if math.IsNaN(float64(nextMaxVal)) {
				panic("training diverged")
			}
			td := r
			if p.Client.PlayMode() != rcsscommon.PlayModeTimeOver {
				td += nextMaxVal
			}
			currentQ := qValues.Get(action)
			currentQVal := currentQ.(float32)
			qValues.Set(action, currentQVal+alpha*(gamma*td-currentQVal))

			err = qLearning.Update(state, qValues)
			if err != nil {
				p.Client.Log(err)
			}

			// S <- S'
			state = nextState
			// time.Sleep(1000 * time.Millisecond)
		}
		epsilon = epsilon * epsilonDecay
		if naiveGames > 0 {
			naiveGames--
		}
		gameCounter++
		timeSinceStart := time.Now().Sub(trainingStart)
		log.Printf("game: %d | return: %f | total time: %s | time/game: %.1f\n", gameCounter, returnValue, timeSinceStart, timeSinceStart.Seconds()/float64(gameCounter))

		// Write return at the end of episode
		returnValues = append(returnValues, returnValue)

		if gameCounter%10 == 0 {
			err = qLearning.Save(weightsFile)
			if err != nil {
				log.Println(err)
			}
			log.Printf("weights saved after %d games\n", gameCounter)
			if gameCounter%50 == 0 {
				file, err := os.Create(returnsFile)
				if err != nil {
					log.Println(err)
				}

				enc := gob.NewEncoder(file)
				err = enc.Encode(returnValues)
				if err != nil {
					log.Println(err)
				}

				file.Close()
				log.Printf("return history saved after %d games\n", gameCounter)
			}
		}
		time.Sleep(2 * time.Second)
	}
}
