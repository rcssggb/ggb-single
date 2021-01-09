package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/rcssggb/ggb-lib/trainerclient"
	"github.com/rcssggb/ggb-single/player"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	hostName := "rcssserver"

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

		serverParams := p.Client.ServerParams()
		var xErr, yErr, tErr float64
		var bXErr, bYErr float64
		var nErr float64

		var estXpos, estYpos, estTpos []float64
		var Xpos, Ypos, Tpos []float64

		var estBallX, estBallY []float64
		var ballXpos, ballYpos []float64

		t.Start()

		for {
			sight := p.Client.See()
			body := p.Client.SenseBody()
			currentTime := p.Client.Time()

			if currentTime == 0 {
				p.Client.Move(-30, 0)
			} else {
				if sight.Ball == nil {
					p.Client.Turn(30)
				} else {
					ballAngle := sight.Ball.Direction + body.HeadAngle
					ballDist := sight.Ball.Distance
					if ballDist < 0.7 {
						p.Client.Kick(20, 0)
					} else {
						p.Client.Dash(50, ballAngle)
						p.Client.TurnNeck(sight.Ball.Direction / 2)
					}
				}
			}
			pAbsPos := t.GlobalPositions().Teams["single-agent"][1]
			bAbsPos := t.GlobalPositions().Ball
			// t.Log(fmt.Sprintf("abs %.2f %.2f %.2f", pAbsPos.X, pAbsPos.Y, pAbsPos.BodyAngle))

			pEstPos := p.GetSelfPos()
			bEstPos := p.GetBall()

			// t.Log(fmt.Sprintf("est %.2f %.2f %.2f", xEstimate, yEstimate, tEstimate))
			nErr++
			xErr = ((nErr-1)/nErr)*xErr + (1/nErr)*math.Abs(pEstPos.X-pAbsPos.X)
			yErr = ((nErr-1)/nErr)*yErr + (1/nErr)*math.Abs(pEstPos.Y-pAbsPos.Y)
			tErr = ((nErr-1)/nErr)*tErr + (1/nErr)*math.Abs(pEstPos.T-pAbsPos.BodyAngle)

			bXErr = ((nErr-1)/nErr)*bXErr + (1/nErr)*math.Abs(bEstPos.X-bAbsPos.X)
			bYErr = ((nErr-1)/nErr)*bYErr + (1/nErr)*math.Abs(bEstPos.Y-bAbsPos.Y)

			// Self position
			estXpos = append(estXpos, pEstPos.X)
			estYpos = append(estYpos, pEstPos.Y)
			estTpos = append(estTpos, pEstPos.T)

			Xpos = append(Xpos, pAbsPos.X)
			Ypos = append(Ypos, pAbsPos.Y)
			Tpos = append(Tpos, pAbsPos.BodyAngle)

			// Ball position
			if bEstPos.NotSeenFor == 0 {
				estBallX = append(estBallX, bEstPos.X)
				estBallY = append(estBallY, bEstPos.Y)

				ballXpos = append(ballXpos, bAbsPos.X)
				ballYpos = append(ballYpos, bAbsPos.Y)
			} else {
				ballXpos = append(ballXpos, bAbsPos.X)
				ballYpos = append(ballYpos, bAbsPos.Y)
			}

			if p.Client.PlayMode() == "time_over" {
				p.Client.Bye()
				break
			}

			err := p.Client.Error()
			for err != nil {
				p.Client.Log(err)
				err = p.Client.Error()
			}

			if currentTime != 0 {
				// time.Sleep(10 * time.Millisecond)
			}

			if serverParams.SynchMode {
				p.Client.DoneSynch()
				t.DoneSynch()
				p.Client.WaitSynch()
				t.WaitSynch()
			} else {
				p.Client.WaitNextStep(currentTime)
				t.WaitNextStep(currentTime)
			}
		}

		t.Log(fmt.Sprintf("Average X Error: %.3f", xErr))
		t.Log(fmt.Sprintf("Average Y Error: %.3f", yErr))
		t.Log(fmt.Sprintf("Average T Error: %.3f", tErr))

		t.Log(fmt.Sprintf("Average Ball X Error: %.3f", bXErr))
		t.Log(fmt.Sprintf("Average Ball Y Error: %.3f", bYErr))

		estPoints := [][]float64{estXpos, estYpos}
		absPoints := [][]float64{Xpos, Ypos}

		ballEstPoints := [][]float64{estBallX, estBallY}
		ballAbsPoints := [][]float64{ballXpos, ballYpos}

		fmt.Println("Saving estimations...")

		estJSON, _ := json.Marshal(estPoints)
		absJSON, _ := json.Marshal(absPoints)
		estTJSON, _ := json.Marshal(estTpos)
		absTJSON, _ := json.Marshal(Tpos)

		ballEstJSON, _ := json.Marshal(ballEstPoints)
		ballAbsJSON, _ := json.Marshal(ballAbsPoints)

		ioutil.WriteFile("data/estJSON.json", estJSON, 0644)
		ioutil.WriteFile("data/absJSON.json", absJSON, 0644)
		ioutil.WriteFile("data/estTJSON.json", estTJSON, 0644)
		ioutil.WriteFile("data/absTJSON.json", absTJSON, 0644)

		ioutil.WriteFile("data/ballEstJSON.json", ballEstJSON, 0644)
		ioutil.WriteFile("data/ballAbsJSON.json", ballAbsJSON, 0644)

		time.Sleep(2 * time.Second)
	}
}
