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
	epsilon := 0.9
	const alpha = 0.1
	const gamma = 0.99
	const epsilonDecay = 0.99996
	const alphaDecay = 0.99999
	const nStates = 282240
	const nActions = 8
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
		if err != nil {
			panic(err)
		}
	}
	qLearning.Gamma = gamma
	qLearning.Alpha = alpha

	_, err = os.Stat(returnsFile)
	returnValues := []float64{}
	if os.IsNotExist(err) {
		log.Println("creating new returns file")
	} else {
		log.Printf("loading return history from %s\n", returnsFile)
		f, err := os.Open(returnsFile)
		if err != nil {
			panic(err)
		}

		dec := gob.NewDecoder(f)
		err = dec.Decode(&returnValues)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
	trainingStart := time.Now()
	lastEnd := trainingStart
	errCount := 0
	for {
		p, err := player.NewPlayer("single-agent", hostName)
		if err != nil {
			errCount++
			if errCount > 10 {
				panic(err)
			}
			log.Println(err)
			continue
		}

		t, err := trainerclient.NewTrainerClient(hostName)
		if err != nil {
			errCount++
			if errCount > 10 {
				panic(err)
			}
			log.Println(err)
			continue
		}

		t.EarOn()
		t.EyeOn()
		p.Client.SynchSee()
		p.Client.ChangeView(rcsscommon.ViewWidthNarrow, rcsscommon.ViewQualityHigh)

		time.Sleep(10 * time.Millisecond)

		// Initialize S
		state := p.State()
		startX, startY := rcsscommon.RandomPosition()
		if startX > 0 {
			startX = -startX
		}
		startT := rand.Float64()*360 - 180
		t.MovePlayer("single-agent", 1, startX, startY, startT, 0, 0)
		t.Start()
		lastGoalTime := -1
		currentTime := 0
		returnValue := float64(0)
		for {
			if p.Client.PlayMode() == rcsscommon.PlayModeTimeOver {
				p.Client.Log(p.Client.Bye())
				break
			}
			if p.Client.PlayMode() == rcsscommon.PlayModeBeforeKickOff {
				t.MovePlayer("single-agent", 1, startX, startY, startT, 0, 0)
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

			nextState := p.State()
			r := float64(0)

			ppos := t.GlobalPositions().Teams["single-agent"][1]
			bpos := t.GlobalPositions().Ball

			distToBall := math.Sqrt(math.Pow(bpos.X-ppos.X, 2) + math.Pow(bpos.Y-ppos.Y, 2))
			r += -distToBall * 0.001 / 6000.0

			r += bpos.DeltaX / 6000.0

			// gx, gy := rcsscommon.FlagRightGoal.Position()
			// r += -math.Sqrt(math.Pow(bpos.X-gx, 2)+math.Pow(bpos.Y-gy, 2)) / gx * 0.0001

			if currentTime > lastGoalTime {
				lastGoalTime = currentTime
				if p.Client.PlayMode() == rcsscommon.PlayModeGoalL {
					r += 1.0
				} else if p.Client.PlayMode() == rcsscommon.PlayModeGoalR {
					r += -1.0
				}
			}

			returnValue += r

			// Update Q towards target
			qLearning.Update(state, action, r, nextState)

			// S <- S'
			state = nextState
			// time.Sleep(10 * time.Microsecond)
		}
		if naiveGames > 0 {
			naiveGames--
		}
		gameCounter++
		epsilon *= epsilonDecay
		qLearning.Alpha *= alphaDecay

		now := time.Now()
		timeSinceStart := now.Sub(trainingStart)
		gameTime := now.Sub(lastEnd)
		lastEnd = now
		log.Printf("game: %d | return: %f | game time: %.3fs | avg time: %.2fs\n",
			gameCounter,
			returnValue,
			gameTime.Seconds(),
			timeSinceStart.Seconds()/float64(gameCounter))

		// Write return at the end of episode
		returnValues = append(returnValues, returnValue)

		if gameCounter%100 == 0 {
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
			log.Printf("training time = %s\n", timeSinceStart)
			log.Printf("alpha = %f\n", qLearning.Alpha)
			log.Printf("epsilon = %f\n", epsilon)

		}
		time.Sleep(1400 * time.Millisecond)
	}
}
