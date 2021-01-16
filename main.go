package main

import (
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

		t.Start()
		for {
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
						if math.Abs(ball.Y) < 15 {
							if ball.X > 0 {
								p.Client.Kick(20, 180-body.T)
							} else {
								p.Client.Kick(20, -body.T)
							}
						} else {
							if ball.Y > 0 {
								p.Client.Kick(20, -90-body.T)
							} else {
								p.Client.Kick(20, 90-body.T)
							}
						}
					} else {
						p.Client.Dash(30, ballAngle)
						p.Client.TurnNeck(ball.Direction / 2)
					}
				}
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

			t.DoneSynch()
			p.WaitCycle()
		}

		time.Sleep(5 * time.Second)
	}
}
