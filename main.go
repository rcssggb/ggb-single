package main

import (
	"log"
	"time"

	"github.com/rcssggb/ggb-lib/playerclient"
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

		time.Sleep(2 * time.Second)

		serverParams := p.ServerParams()
		for {
			sight := p.See()
			body := p.SenseBody()
			currentTime := p.Time()

			if sight.Ball == nil {
				p.Turn(30)
			} else {
				ballAngle := sight.Ball.Direction + body.HeadAngle
				ballDist := sight.Ball.Distance
				if ballDist < 0.7 {
					p.Kick(20, 0)
				} else {
					p.Dash(50, ballAngle)
					p.TurnNeck(sight.Ball.Direction)
				}
			}

			if p.PlayMode() == "time_over" {
				p.Bye()
				break
			}

			err := p.Error()
			for err != nil {
				p.Log(err)
				err = p.Error()
			}

			if serverParams.SynchMode {
				p.DoneSynch()
				p.WaitSynch()
			} else {
				p.WaitNextStep(currentTime)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
