package main

import (
	"log"
	"math"
	"os"
	"path"
	"time"

	"github.com/rcssggb/ggb-lib/rcsscommon"
	"github.com/rcssggb/ggb-lib/trainerclient"
	"github.com/rcssggb/ggb-single/player"
)

func main() {
	epsilon := 0.2
	const alpha = float32(1)
	const epsilonDecay = 0.9999
	naiveGames := 0
	gameCounter := 0
	// weightsFileA := "weightsA.rln"
	// weightsFileB := "weightsB.rln"
	// returnsFile := "./data/returns.rln"

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

	// var qLearningA, qLearningB *q.QLearning

	// _, err = os.Stat(weightsFileA)
	// if os.IsNotExist(err) {
	// 	log.Println("creating new agent")
	// 	qLearningA, err = q.Init()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// } else {
	// 	log.Printf("loading agent from %s\n", weightsFileA)
	// 	qLearningA, err = q.Load(weightsFileA)
	// }

	// _, err = os.Stat(weightsFileB)
	// if os.IsNotExist(err) {
	// 	log.Println("creating new agent")
	// 	qLearningB, err = q.Init()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// } else {
	// 	log.Printf("loading agent from %s\n", weightsFileB)
	// 	qLearningB, err = q.Load(weightsFileB)
	// }

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
		// state := q.Slice2Tensor(p.State())
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

			if (currentTime % 1000) == 0 {
				x, y := rcsscommon.RandomPosition()
				t.MoveBall(x, y, 0, 0)
			}

			// Choose A from S using policy derived from Q (e.g., epsilon-greedy)
			// qValuesA, err := qLearningA.ActionValues(state)
			// if err != nil {
			// 	p.Client.Log(err)
			// }
			// qValuesB, err := qLearningB.ActionValues(state)
			// if err != nil {
			// 	p.Client.Log(err)
			// }

			// var qValues *tensor.Dense
			// qSum, err := qValuesA.Add(qValuesB)
			// if err != nil {
			// 	p.Client.Log(err)
			// }
			// qValues, err = qSum.DivScalar(float32(2.0), true)
			// if err != nil {
			// 	p.Client.Log(err)
			// }

			// var action int
			// takeRandomAction := rand.Float64() < epsilon
			// if takeRandomAction {
			// 	action = rand.Intn(16)
			// } else {
			// 	if naiveGames > 0 {
			// 		action = p.NaivePolicy()
			// 	} else {
			// 		maxActionTensor, err := qValues.Argmax(1)
			// 		if err != nil {
			// 			p.Client.Log(err)
			// 		}
			// 		action = maxActionTensor.Data().([]int)[0]
			// 	}
			// }

			// Try and arrange behaviors for testing
			var action int
			if p.GetBall().NotSeenFor == 0 {
				if p.GetSelfData().X > 30 {
					// Shoot ball
					action = 2
				} else {
					// Lead ball
					action = 1
				}
			} else {
				// fmt.Println("ball not seen for ", p.GetBall().NotSeenFor)
				// Locate ball
				action = 0
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
			// nextState := q.Slice2Tensor(p.State())
			currentTime = p.Client.Time()
			r := float32(0)

			r = float32(math.Abs(t.GlobalPositions().Teams["single-agent"][1].BodyAngle)/90.0 - 1.0)

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

			// // Update Q towards target
			// nextActionValuesA, err := qLearningA.ActionValues(state)
			// if err != nil {
			// 	p.Client.Log(err)
			// }
			// nextActionValuesB, err := qLearningB.ActionValues(state)
			// if err != nil {
			// 	p.Client.Log(err)
			// }
			// if rand.Float32() < 0.5 {
			// 	maxActionCoord, err := nextActionValuesA.Argmax(1)
			// 	if err != nil {
			// 		p.Client.Log(err)
			// 	}
			// 	maxActionCoordVal := maxActionCoord.Data().([]int)[0]
			// 	nextMax := nextActionValuesB.Get(maxActionCoordVal)
			// 	if err != nil {
			// 		p.Client.Log(err)
			// 	}
			// 	nextMaxVal := nextMax.(float32)
			// 	if math.IsNaN(float64(nextMaxVal)) {
			// 		panic("training diverged")
			// 	}
			// 	td := r
			// 	if p.Client.PlayMode() != rcsscommon.PlayModeTimeOver {
			// 		td += nextMaxVal
			// 	}
			// 	currentQ := qValuesA.Get(action)
			// 	currentQVal := currentQ.(float32)
			// 	qValuesA.Set(action, currentQVal+alpha*(td-currentQVal))
			// 	err = qLearningA.UpdateWithBatch(state, qValuesA)
			// 	if err != nil {
			// 		p.Client.Log(err)
			// 	}
			// } else {
			// 	maxActionCoord, err := nextActionValuesB.Argmax(1)
			// 	if err != nil {
			// 		p.Client.Log(err)
			// 	}
			// 	maxActionCoordVal := maxActionCoord.Data().([]int)[0]
			// 	nextMax := nextActionValuesA.Get(maxActionCoordVal)
			// 	if err != nil {
			// 		p.Client.Log(err)
			// 	}
			// 	nextMaxVal := nextMax.(float32)
			// 	if math.IsNaN(float64(nextMaxVal)) {
			// 		panic("training diverged")
			// 	}
			// 	td := r
			// 	if p.Client.PlayMode() != rcsscommon.PlayModeTimeOver {
			// 		td += nextMaxVal
			// 	}
			// 	currentQ := qValuesB.Get(action)
			// 	currentQVal := currentQ.(float32)
			// 	qValuesB.Set(action, currentQVal+alpha*(td-currentQVal))
			// 	err = qLearningB.UpdateWithBatch(state, qValuesB)
			// 	if err != nil {
			// 		p.Client.Log(err)
			// 	}
			// }

			// S <- S'
			// state = nextState
		}
		if naiveGames > 0 {
			naiveGames--
		}
		gameCounter++
		timeSinceStart := time.Now().Sub(trainingStart)
		log.Printf("game: %d | return: %f | total time: %s | time/game: %.1f\n", gameCounter, returnValue, timeSinceStart, timeSinceStart.Seconds()/float64(gameCounter))

		// Write return at the end of episode
		returnValues = append(returnValues, returnValue)

		// if gameCounter%10 == 0 {
		// 	err = qLearningA.Save(weightsFileA)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	err = qLearningB.Save(weightsFileB)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	log.Printf("weights saved after %d games\n", gameCounter)
		// 	if gameCounter%50 == 0 {
		// 		file, err := os.Create(returnsFile)
		// 		if err != nil {
		// 			log.Println(err)
		// 		}

		// 		enc := gob.NewEncoder(file)
		// 		err = enc.Encode(returnValues)
		// 		if err != nil {
		// 			log.Println(err)
		// 		}

		// 		file.Close()
		// 		log.Printf("return history saved after %d games\n", gameCounter)
		// 	}
		// }
		time.Sleep(2 * time.Second)
	}
}
