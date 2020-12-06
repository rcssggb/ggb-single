package main

import (
	"log"
	"time"

	"github.com/rcssggb/ggb-lib/playerclient"
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
			body := p.SenseBody()
			currentTime := p.Time()

			if currentTime == 0 {
				p.Move(-30, 0)
			} else {
				if sight.Ball == nil {
					p.Turn(30)
				} else {
					ballAngle := sight.Ball.Direction + body.HeadAngle
					ballDist := sight.Ball.Distance
					if ballDist < 0.7 {
						p.Kick(20, 0)
					} else {
						p.Dash(50, ballAngle)
						p.TurnNeck(sight.Ball.Direction / 2)
					}
				}
			}
			pAbsPos := t.GlobalPositions().Teams["single-agent"][1]
			t.Log(fmt.Sprintf("%f, %f, %f", pAbsPos.X, pAbsPos.Y, pAbsPos.BodyAngle))

			if p.PlayMode() == "time_over" {
				p.Bye()
				break
			}

			err := p.Error()
			for err != nil {
				p.Log(err)
				err = p.Error()
			}

			time.Sleep(500 * time.Millisecond)

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
