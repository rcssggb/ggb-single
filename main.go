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
	epsilon := 0.5
	const alpha = float32(1)
	const epsilonDecay = 0.999
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

	log.Printf("starting training with // epsilon = %f // alpha = %f // epsilonDecay = %f // naiveGames = %d", epsilon, alpha, epsilonDecay, naiveGames)

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
		returnValue := float64(0)
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
				action = rand.Intn(4)
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
			nextState := q.Slice2Tensor(p.State())
			currentTime = p.Client.Time()
			// r := float32(0)

			// // r = float32(math.Abs(t.GlobalPositions().Teams["single-agent"][1].BodyAngle)/90.0 - 1.0)

			// // if ball.NotSeenFor == 0 {
			// // 	ballDist := float32(ball.Distance)
			// // 	if ballDist < 0.7 {
			// // 		ballDist = 0.7
			// // 		epsilon *= epsilonDecay
			// // 	}
			// // 	r = 1.0 / ballDist
			// // }

			// if p.Client.PlayMode() == rcsscommon.PlayModeGoalL && currentTime > lastGoalTime {
			// 	lastGoalTime = currentTime
			// 	r = 1
			// 	p.Client.Log("goal!")
			// }

			// if p.Client.PlayMode() == rcsscommon.PlayModeGoalR && currentTime > lastGoalTime {
			// 	lastGoalTime = currentTime
			// 	r = -1
			// 	p.Client.Log("goal against, bad!")
			// }

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
			qValues.Set(action, currentQVal+alpha*(td-currentQVal))
			err = qLearning.UpdateWithBatch(state, qValues)
			if err != nil {
				p.Client.Log(err)
			}

			// S <- S'
			state = nextState
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
