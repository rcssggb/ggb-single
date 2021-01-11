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

		var xErr, yErr, tErr float64
		var xVErr, yVErr float64
		var bXErr, bYErr float64
		var bVXErr, bVYErr float64
		var nErr float64

		var estXpos, estYpos, estTpos []float64
		var estXVel, estYVel []float64
		var Xpos, Ypos, Tpos []float64
		var XVel, YVel []float64

		var estBallX, estBallY []float64
		var ballXpos, ballYpos []float64

		t.Start()

		for {
			// sight := p.Client.See()
			// body := p.Client.SenseBody()
			currentTime := p.Client.Time()

			if currentTime != 0 {
				// time.Sleep(200 * time.Millisecond)
			}

			ball := p.GetBall()
			body := p.GetSelfData()
			if currentTime == 0 {
				p.Client.Move(-30, 0)
			} else {
				if ball.NotSeenFor != 0 {
					p.Client.Turn(20)
				} else {
					ballAngle := ball.Direction + body.NeckAngle
					ballDist := ball.Distance
					if ballDist < 0.7 {
						if math.Abs(ball.Y) > 15 {
							if ball.X > 0 {
								p.Client.Kick(40, 180-body.T)
							} else {
								p.Client.Kick(40, -body.T)
							}
						} else {
							if ball.Y > 0 {
								p.Client.Kick(40, -90-body.T)
							} else {
								p.Client.Kick(40, 90-body.T)
							}
						}
					} else {
						p.Client.Dash(30, ballAngle)
						p.Client.TurnNeck(ball.Direction / 2)
					}
				}
			}

			t.MovePlayer("single-agent", 2, ball.X+2, ball.Y+2, 0, 0, 0)

			pAbsPos := t.GlobalPositions().Teams["single-agent"][1]
			bAbsPos := t.GlobalPositions().Ball
			// t.Log(fmt.Sprintf("abs %.2f %.2f %.2f", pAbsPos.X, pAbsPos.Y, pAbsPos.BodyAngle))

			pEstPos := body
			bEstPos := ball

			// t.Log(fmt.Sprintf("est %.2f %.2f %.2f", xEstimate, yEstimate, tEstimate))
			nErr++
			xErr = ((nErr-1)/nErr)*xErr + (1/nErr)*math.Abs(pEstPos.X-pAbsPos.X)
			yErr = ((nErr-1)/nErr)*yErr + (1/nErr)*math.Abs(pEstPos.Y-pAbsPos.Y)
			absTErr := math.Abs(pEstPos.T - pAbsPos.BodyAngle)
			if absTErr > 180 {
				absTErr = 360 - absTErr
			}
			tErr = ((nErr-1)/nErr)*tErr + (1/nErr)*absTErr

			xVErr = ((nErr-1)/nErr)*xVErr + (1/nErr)*math.Abs(pEstPos.VelX-pAbsPos.DeltaX)
			yVErr = ((nErr-1)/nErr)*yVErr + (1/nErr)*math.Abs(pEstPos.VelY-pAbsPos.DeltaY)

			// t.Log(fmt.Sprintf("Estimated X: %f, Absolute X: %f\n", pEstPos.VelX, pAbsPos.DeltaX))
			// t.Log(fmt.Sprintf("Estimated Y: %f, Absolute Y: %f\n", pEstPos.VelY, pAbsPos.DeltaY))

			bXErr = ((nErr-1)/nErr)*bXErr + (1/nErr)*math.Abs(bEstPos.X-bAbsPos.X)
			bYErr = ((nErr-1)/nErr)*bYErr + (1/nErr)*math.Abs(bEstPos.Y-bAbsPos.Y)
			bVXErr = ((nErr-1)/nErr)*bVXErr + (1/nErr)*math.Abs(bEstPos.VelX-bAbsPos.DeltaX)
			bVYErr = ((nErr-1)/nErr)*bVYErr + (1/nErr)*math.Abs(bEstPos.VelY-bAbsPos.DeltaY)

			// Self position
			estXpos = append(estXpos, pEstPos.X)
			estYpos = append(estYpos, pEstPos.Y)
			estTpos = append(estTpos, pEstPos.T)

			Xpos = append(Xpos, pAbsPos.X)
			Ypos = append(Ypos, pAbsPos.Y)
			Tpos = append(Tpos, pAbsPos.BodyAngle)

			// Self Velocity
			estXVel = append(estXVel, pEstPos.VelX)
			estYVel = append(estYVel, pEstPos.VelY)

			XVel = append(XVel, pAbsPos.DeltaX)
			YVel = append(YVel, pAbsPos.DeltaY)

			// Ball position
			estBallX = append(estBallX, bEstPos.X)
			estBallY = append(estBallY, bEstPos.Y)

			ballXpos = append(ballXpos, bAbsPos.X)
			ballYpos = append(ballYpos, bAbsPos.Y)

			if p.Client.PlayMode() == "time_over" {
				p.Client.Bye()
				break
			}

			err := p.Client.Error()
			for err != nil {
				p.Client.Log(err)
				err = p.Client.Error()
			}

			t.DoneSynch()
			p.WaitCycle()
		}

		t.Log(fmt.Sprintf("Average X Error: %.3f", xErr))
		t.Log(fmt.Sprintf("Average Y Error: %.3f", yErr))
		t.Log(fmt.Sprintf("Average T Error: %.3f", tErr))

		t.Log(fmt.Sprintf("Average VelX Error: %.3f", xVErr))
		t.Log(fmt.Sprintf("Average VelY Error: %.3f", yVErr))

		t.Log(fmt.Sprintf("Average Ball X Error: %.3f", bXErr))
		t.Log(fmt.Sprintf("Average Ball Y Error: %.3f", bYErr))

		t.Log(fmt.Sprintf("Average Ball VelX Error: %.3f", bVXErr))
		t.Log(fmt.Sprintf("Average Ball VelY Error: %.3f", bVYErr))

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

		estXVelJSON, _ := json.Marshal(estXVel)
		estYVelJSON, _ := json.Marshal(estYVel)
		absXVelJSON, _ := json.Marshal(XVel)
		absYVelJSON, _ := json.Marshal(YVel)

		ioutil.WriteFile("data/estJSON.json", estJSON, 0644)
		ioutil.WriteFile("data/absJSON.json", absJSON, 0644)
		ioutil.WriteFile("data/estTJSON.json", estTJSON, 0644)
		ioutil.WriteFile("data/absTJSON.json", absTJSON, 0644)

		ioutil.WriteFile("data/ballEstJSON.json", ballEstJSON, 0644)
		ioutil.WriteFile("data/ballAbsJSON.json", ballAbsJSON, 0644)

		ioutil.WriteFile("data/estXVelJSON.json", estXVelJSON, 0644)
		ioutil.WriteFile("data/estYVelJSON.json", estYVelJSON, 0644)
		ioutil.WriteFile("data/absXVelJSON.json", absXVelJSON, 0644)
		ioutil.WriteFile("data/absYVelJSON.json", absYVelJSON, 0644)

		time.Sleep(2 * time.Second)
	}
}
