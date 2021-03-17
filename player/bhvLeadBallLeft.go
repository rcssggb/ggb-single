package player

import (
	"math"
)

func (p *Player) bhvLeadBallLeft() string {
	cmd := ""
	body := p.GetSelfData()
	ball := p.GetBall()
	if ball.Distance < 0.8 {
		x := body.X - 10
		y := body.Y
		angle := (180.0/math.Pi)*math.Atan2(y-body.Y, x-body.X) - body.T
		cmd += p.Client.Kick(5, angle)
	}
	return cmd
}
