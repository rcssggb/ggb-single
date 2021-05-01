package main

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
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

type trainingInfo struct {
	Epsilon   float64
	Alpha     float32
	GameCount int
}

func main() {
	epsilon := 0.9
	alpha := float32(0.1)
	const gamma = 0.99
	const epsilonDecay = 0.9999
	const alphaDecay = 0.99999
	naiveGames := 0
	gameCounter := 0
	saveEvery := 5
	weightsFile := "weights.rln"
	returnsFile := "./data/returns.rln"
	infoFile := "info.json"

	logName := time.Now().String() + ".log"
	file, err := os.OpenFile(path.Join("logs", logName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var info trainingInfo
	_, err = os.Stat(infoFile)
	if !os.IsNotExist(err) {
		i, _ := ioutil.ReadFile(infoFile)
		err = json.Unmarshal(i, &info)
		if err != nil {
			log.Fatal(err)
		}
		alpha = info.Alpha
		epsilon = info.Epsilon
		gameCounter = info.GameCount
	} else {
		info.Alpha = alpha
		info.Epsilon = epsilon
		info.GameCount = gameCounter
	}

	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	log.Printf("starting training with\n alpha = %f\n alphaDecay = %f\n epsilon = %f\n epsilonDecay = %f\n naiveGames = %d", alpha, alphaDecay, epsilon, epsilonDecay, naiveGames)

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

		// Choose A from S using policy derived from Q (e.g., epsilon-greedy)
		qValues, err := qLearning.ActionValues(state)
		if err != nil {
			p.Client.Log(err)
		}

		var action int
		takeRandomAction := rand.Float64() < epsilon
		if takeRandomAction {
			action = rand.Intn(qLearning.NBehaviors)
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

		for {
			if p.Client.PlayMode() == rcsscommon.PlayModeTimeOver {
				p.Client.Log(p.Client.Bye())
				break
			}

			// Take action A
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

			// Observe S'
			nextState := q.Slice2Tensor(p.State())
			currentTime = p.Client.Time()
			r := float32(0)

			ppos := t.GlobalPositions().Teams["single-agent"][1]
			bpos := t.GlobalPositions().Ball

			distToBall := math.Sqrt(math.Pow(bpos.X-ppos.X, 2) + math.Pow(bpos.Y-ppos.Y, 2))
			r += float32(-distToBall) * 0.001 / 6000.0

			r += float32(bpos.DeltaX) / 6000.0

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

			// Choose A' from S' using policy derived from Q (e.g., epsilon-greedy)
			nextQValues, err := qLearning.ActionValues(nextState)
			if err != nil {
				p.Client.Log(err)
			}
			var nextAction int
			takeRandomAction := rand.Float64() < epsilon
			if takeRandomAction {
				nextAction = rand.Intn(qLearning.NBehaviors)
			} else {
				if naiveGames > 0 {
					nextAction = p.NaiveBehaviorPolicy()
				} else {
					maxActionTensor, err := nextQValues.Argmax(1)
					if err != nil {
						p.Client.Log(err)
					}
					nextAction = maxActionTensor.Data().([]int)[0]
				}
			}

			// Check if training diverged
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

			// Update Q towards target
			currentQ := qValues.Get(action)
			currentQVal := currentQ.(float32)

			nextQ := qValues.Get(nextAction)
			nextQVal := nextQ.(float32)

			td := r
			if p.Client.PlayMode() != rcsscommon.PlayModeTimeOver {
				td += gamma * nextQVal
			}

			qValues.Set(action, currentQVal+alpha*(td-currentQVal))

			err = qLearning.Update(state, qValues)
			if err != nil {
				p.Client.Log(err)
			}

			// S <- S'
			// A <- A'
			state = nextState
			action = nextAction
		}
		epsilon = epsilon * epsilonDecay
		alpha = alpha * alphaDecay
		if naiveGames > 0 {
			naiveGames--
		}
		gameCounter++
		timeSinceStart := time.Now().Sub(trainingStart)
		log.Printf("game: %d | return: %f | total time: %s | time/game: %.1f\n", gameCounter, returnValue, timeSinceStart, timeSinceStart.Seconds()/float64(gameCounter))

		// Write return at the end of episode
		returnValues = append(returnValues, returnValue)

		if gameCounter%saveEvery == 0 {
			err = qLearning.Save(weightsFile)
			if err != nil {
				log.Println(err)
			}
			log.Printf("current epsilon = %f\n", epsilon)
			log.Printf("weights saved after %d games\n", gameCounter)
			if gameCounter%saveEvery == 0 {
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

				info.Alpha = alpha
				info.Epsilon = epsilon
				info.GameCount = gameCounter
				i, err := json.Marshal(info)
				if err != nil {
					log.Println(err)
				}
				err = ioutil.WriteFile("info.json", i, 0666)
				if err != nil {
					log.Println(err)
				}

				log.Printf("return history saved after %d games\n", gameCounter)
				log.Printf("training time = %s\n", timeSinceStart)
				log.Printf("alpha = %f\n", alpha)
				log.Printf("epsilon = %f\n", epsilon)
			}
		}
		time.Sleep(2 * time.Second)
	}
}
