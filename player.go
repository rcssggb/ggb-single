package main

import "github.com/rcssggb/ggb-lib/playerclient"

func player(c *playerclient.Client) {
	serverParams := c.ServerParams()
	for {
		sight := c.See()
		body := c.SenseBody()
		currentTime := c.Time()

		if sight.Ball == nil {
			c.Turn(30)
		} else {
			ballAngle := sight.Ball.Direction + body.HeadAngle
			ballDist := sight.Ball.Distance
			if ballDist < 0.7 {
				c.Kick(20, 0)
			} else {
				c.Dash(50, ballAngle)
				c.TurnNeck(sight.Ball.Direction)
			}
		}

		if c.PlayMode() == "time_over" {
			c.Bye()
			break
		}

		err := c.Error()
		for err != nil {
			c.Log(err)
			err = c.Error()
		}

		if serverParams.SynchMode {
			c.DoneSynch()
			c.WaitSynch()
		} else {
			c.WaitNextStep(currentTime)
		}
	}
}
