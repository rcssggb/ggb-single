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
		var nErr float64
		estXpos := []float64{}
		estYpos := []float64{}
		estTpos := []float64{}
		Xpos := []float64{}
		Ypos := []float64{}
		Tpos := []float64{}
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
			// t.Log(fmt.Sprintf("abs %.2f %.2f %.2f", pAbsPos.X, pAbsPos.Y, pAbsPos.BodyAngle))

			pEstPos := p.GetSelfPos()

			// t.Log(fmt.Sprintf("est %.2f %.2f %.2f", xEstimate, yEstimate, tEstimate))
			nErr++
			xErr = ((nErr-1)/nErr)*xErr + (1/nErr)*math.Abs(pEstPos.X-pAbsPos.X)
			yErr = ((nErr-1)/nErr)*yErr + (1/nErr)*math.Abs(pEstPos.Y-pAbsPos.Y)
			tErr = ((nErr-1)/nErr)*tErr + (1/nErr)*math.Abs(pEstPos.T-pAbsPos.BodyAngle)

			estXpos = append(estXpos, pEstPos.X)
			estYpos = append(estYpos, pEstPos.Y)
			estTpos = append(estTpos, pEstPos.T)
			Xpos = append(Xpos, pAbsPos.X)
			Ypos = append(Ypos, pAbsPos.Y)
			Tpos = append(Tpos, pAbsPos.BodyAngle)

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

		estPoints := [][]float64{estXpos, estYpos}
		absPoints := [][]float64{Xpos, Ypos}

		fmt.Println("Saving estimations...")

		estJSON, _ := json.Marshal(estPoints)
		absJSON, _ := json.Marshal(absPoints)
		estTJSON, _ := json.Marshal(estTpos)
		absTJSON, _ := json.Marshal(Tpos)

		ioutil.WriteFile("data/estJSON.json", estJSON, 0644)
		ioutil.WriteFile("data/absJSON.json", absJSON, 0644)
		ioutil.WriteFile("data/estTJSON.json", estTJSON, 0644)
		ioutil.WriteFile("data/absTJSON.json", absTJSON, 0644)

		time.Sleep(2 * time.Second)
	}
}
