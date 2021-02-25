package player

import (
	"math"
)

func (p *Player) bhvLeadBallRight() string {
	cmd := ""
	body := p.GetSelfData()
	x := body.X + 10
	y := body.Y
	angle := (180.0/math.Pi)*math.Atan2(y-body.Y, x-body.X) - body.T
	cmd += p.Client.Kick(5, angle)
	return cmd
}
