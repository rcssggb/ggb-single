package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rcssggb/ggb-lib/playerclient"
	"github.com/rcssggb/ggb-lib/rcsscommon"
	"github.com/rcssggb/ggb-lib/trainerclient"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	hostName := "rcssserver"

	for {
		p, err := playerclient.NewPlayerClient("single-agent", hostName)
		if err != nil {
			log.Println(err)
			p.Bye()
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

		serverParams := p.ServerParams()
		for {
			sight := p.See()
			// body := p.SenseBody()
			currentTime := p.Time()

			if currentTime == 0 {
				p.Move(-30, 0)
			} else {
				// if sight.Ball == nil {
				// 	p.Turn(30)
				// } else {
				// 	ballAngle := sight.Ball.Direction + body.HeadAngle
				// 	ballDist := sight.Ball.Distance
				// 	if ballDist < 0.7 {
				// 		p.Kick(20, 0)
				// 	} else {
				// 		p.Dash(50, ballAngle)
				// 		p.TurnNeck(sight.Ball.Direction / 2)
				// 	}
				// }
				p.Turn(10)
			}
			pAbsPos := t.GlobalPositions().Teams["single-agent"][1]
			t.Log(fmt.Sprintf("tAbsolute %f", pAbsPos.BodyAngle))

			var tEstimate float64
			if sight.Lines.Len() > 0 {
				closestLine := sight.Lines[0]
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

			if tEstimate > 180 {
				tEstimate -= 360
			} else if tEstimate < -180 {
				tEstimate += 360
			}

			p.Log(fmt.Sprintf("tEstimate %f", tEstimate))

			if p.PlayMode() == "time_over" {
				p.Bye()
				break
			}

			err := p.Error()
			for err != nil {
				p.Log(err)
				err = p.Error()
			}

			if currentTime != 0 {
				time.Sleep(200 * time.Millisecond)
			}

			if serverParams.SynchMode {
				p.DoneSynch()
				t.DoneSynch()
				p.WaitSynch()
				t.WaitSynch()
			} else {
				p.WaitNextStep(currentTime)
				t.WaitNextStep(currentTime)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
