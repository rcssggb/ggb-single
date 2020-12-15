package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/rcssggb/ggb-lib/playerclient/parser"
	"github.com/rcssggb/ggb-lib/rcsscommon"
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

			var tEstimate, xEstimate, yEstimate float64
			var closestLine *parser.LineData
			closestLine = nil
			if sight.Lines.Len() > 0 {
				closestLine = &sight.Lines[0]
				lDir := closestLine.Direction
				if lDir < 0 {
					lDir += 90
				} else {
					lDir -= 90
				}
				switch closestLine.ID {
				case rcsscommon.LineRight:
					tEstimate = 0 - lDir
				case rcsscommon.LineBottom:
					tEstimate = 90 - lDir
				case rcsscommon.LineLeft:
					tEstimate = 180 - lDir
				case rcsscommon.LineTop:
					tEstimate = -90 - lDir
				}
			}

			// If you see 2 lines it means you're outside the field
			if sight.Lines.Len() >= 2 {
				tEstimate += 180
			}

			tEstimate -= body.HeadAngle

			if tEstimate > 180 {
				tEstimate -= 360
			} else if tEstimate < -180 {
				tEstimate += 360
			}

			if sight.Flags.Len() > 0 {
				var xAcc, yAcc float64 = 0, 0
				for _, f := range sight.Flags {
					xFlag, yFlag := f.ID.Position()
					absAngle := (3.14159 / 180.0) * (f.Direction + tEstimate + body.HeadAngle)
					xTmp := xFlag - math.Cos(absAngle)*f.Distance
					yTmp := yFlag - math.Sin(absAngle)*f.Distance
					// p.Client.Log(fmt.Sprintf(
					// 	"ID:%d Dist:%.2f, Dir: %.2f, xEst: %.2f, yEst: %.2f",
					// 	f.ID, f.Distance, absAngle, xTmp, yTmp))
					xAcc += xTmp
					yAcc += yTmp
				}
				xEstimate = xAcc / (float64)(sight.Flags.Len())
				yEstimate = yAcc / (float64)(sight.Flags.Len())
			}

			// t.Log(fmt.Sprintf("est %.2f %.2f %.2f", xEstimate, yEstimate, tEstimate))
			nErr++
			xErr = ((nErr-1)/nErr)*xErr + (1/nErr)*math.Abs(xEstimate-pAbsPos.X)
			yErr = ((nErr-1)/nErr)*yErr + (1/nErr)*math.Abs(yEstimate-pAbsPos.Y)
			tErr = ((nErr-1)/nErr)*tErr + (1/nErr)*math.Abs(tEstimate-pAbsPos.BodyAngle)

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

		time.Sleep(2 * time.Second)
	}
}
