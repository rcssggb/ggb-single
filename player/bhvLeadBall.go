package player

import (
	"math"
)

func (p *Player) bhvLeadBall() string {
	cmd := ""
	ball := p.GetBall()
	body := p.GetSelfData()
	if ball.Distance < 0.7 {
		x := body.X + 10
		y := body.Y
		angle := (180.0/math.Pi)*math.Atan2(y-body.Y, x-body.X) - body.T
		cmd += p.Client.Kick(5, angle)
	}
	// } else {
	// 	ballAngle := ball.Direction + body.NeckAngle
	// 	if -15 < ballAngle && ballAngle < 15 {
	// 		cmd += p.Client.Dash(60, ballAngle)
	// 	} else {
	// 		cmd += p.Client.Turn(ballAngle)
	// 	}
	// }
	return cmd
}
