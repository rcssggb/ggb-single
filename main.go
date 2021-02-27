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
	epsilon := 1.0
	const alpha = 0.5
	const gamma = 0.99
	const epsilonDecay = 0.999
	const alphaDecay = 0.9995
	const nStates = 144
	const nActions = 6
	naiveGames := 0
	gameCounter := 0
	qTableFile := "qtable.rln"
	returnsFile := "./data/returns.rln"

	logName := time.Now().String() + ".log"
	file, err := os.OpenFile(path.Join("logs", logName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	log.Printf("starting training with\n alpha = %f\n alphaDecay = %f\n epsilon = %f\n epsilonDecay = %f\n naiveGames = %d", alpha, alphaDecay, epsilon, epsilonDecay, naiveGames)

	hostName := "rcssserver"

	var qLearning *q.QLearning

	_, err = os.Stat(qTableFile)
	if os.IsNotExist(err) {
		log.Println("creating new agent")
		qLearning = q.Init(nStates, nActions)
	} else {
		log.Printf("loading agent from %s\n", qTableFile)
		qLearning, err = q.Load(qTableFile)
	}
	qLearning.Gamma = gamma
	qLearning.Alpha = alpha

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

		time.Sleep(1 * time.Second)

		// Initialize S
		state := p.State()
		startX, startY := rcsscommon.RandomPosition()
		if startX > 0 {
			startX = -startX
		}
		startT := rand.Float64()*360 - 180
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
				action = rand.Intn(nActions)
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

			ppos := t.GlobalPositions().Teams["single-agent"][1]
			bpos := t.GlobalPositions().Ball
			r = -math.Sqrt(math.Pow(bpos.Y-ppos.Y, 2)+math.Pow(bpos.X-ppos.X, 2)) / 6000.0

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
			err = qLearning.Save(qTableFile)
			if err != nil {
				log.Println(err)
			}
			log.Printf("q table saved after %d games\n", gameCounter)
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
			log.Printf("current parameters\n alpha = %f\n epsilon = %f\n", qLearning.Alpha, epsilon)

		}
		time.Sleep(1500 * time.Millisecond)
	}
}
