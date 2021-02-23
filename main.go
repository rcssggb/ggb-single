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
	q "github.com/rcssggb/ggb-single/tabularq"
)

func main() {
	epsilon := 0.1
	const alpha = 0.2
	const epsilonDecay = 0.999
	const alphaDecay = 1
	naiveGames := 0
	gameCounter := 0
	// weightsFileA := "weightsA.rln"
	// weightsFileB := "weightsB.rln"
	returnsFile := "./data/returns.rln"

	logName := time.Now().String() + ".log"
	file, err := os.OpenFile(path.Join("logs", logName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	log.Printf("starting training with // epsilon = %f // alpha = %f // epsilonDecay = %f // naiveGames = %d", epsilon, alpha, epsilonDecay, naiveGames)

	hostName := "rcssserver"

	var qLearning *q.QLearning

	// _, err = os.Stat(weightsFileA)
	// if os.IsNotExist(err) {
	log.Println("creating new agent")
	qLearning = q.Init(24, 3)
	qLearning.Alpha = alpha
	qLearning.Gamma = 0.99
	// } else {
	// 	log.Printf("loading agent from %s\n", weightsFileA)
	// 	qLearningA, err = q.Load(weightsFileA)
	// }

	returnValues := []float64{}
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
		state := p.State()
		startX, startY := rcsscommon.RandomPosition()
		startT := rand.Float64()*360 - 180
		if startX > 0 {
			startX = -startX
		}
		t.MovePlayer("single-agent", 1, startX, startY, startT, 0, 0)
		t.Start()
		// lastGoalTime := -1
		currentTime := 0
		returnValue := float64(0)
		for {
			if p.Client.PlayMode() == rcsscommon.PlayModeTimeOver {
				p.Client.Log(p.Client.Bye())
				break
			}

			// Choose A from S using policy derived from Q (e.g., epsilon-greedy)
			var action int
			takeRandomAction := rand.Float64() < epsilon
			if takeRandomAction {
				action = rand.Intn(3)
			} else {
				if naiveGames > 0 {
					action = p.NaivePolicy()
				} else {
					action = qLearning.GreedyAction(state)
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
			nextState := p.State()
			currentTime = p.Client.Time()
			r := float64(0)

			r = float64((math.Abs(t.GlobalPositions().Teams["single-agent"][1].BodyAngle) - 90) / 90.0)

			// ball := p.GetBall()
			// if ball.NotSeenFor == 0 {
			// 	ballDist := float32(ball.Distance)
			// 	if ballDist < 0.7 {
			// 		ballDist = 0.7
			// 		epsilon *= epsilonDecay
			// 	}
			// 	r = 1.0 / ballDist
			// }

			// if p.Client.PlayMode() == rcsscommon.PlayModeGoalL && currentTime > lastGoalTime {
			// 	lastGoalTime = currentTime
			// 	r = 1
			// 	p.Client.Log("goal!")
			// 	epsilon = epsilon * epsilonDecay
			// }

			returnValue += r

			// Update Q towards target
			qLearning.Update(state, action, r, nextState)

			// S <- S'
			state = nextState
		}
		if naiveGames > 0 {
			naiveGames--
		}
		gameCounter++
		epsilon *= epsilonDecay
		qLearning.Alpha *= alphaDecay
		timeSinceStart := time.Now().Sub(trainingStart)
		log.Printf("game: %d | return: %f | total time: %s | time/game: %.1f\n", gameCounter, returnValue, timeSinceStart, timeSinceStart.Seconds()/float64(gameCounter))

		// Write return at the end of episode
		returnValues = append(returnValues, returnValue)

		if gameCounter%50 == 0 {
			// err = qLearningA.Save(weightsFileA)
			// if err != nil {
			// 	log.Println(err)
			// }
			// err = qLearningB.Save(weightsFileB)
			// if err != nil {
			// 	log.Println(err)
			// }
			// log.Printf("weights saved after %d games\n", gameCounter)
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
		time.Sleep(1500 * time.Millisecond)
	}
}
