package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/rcsscommon"
)

func (p *Player) bhvShootToGoalR() string {
	cmd := ""
	body := p.GetSelfData()
	ball := p.GetBall()
	if ball.Distance < 0.8 {
		x, y := rcsscommon.FlagRightGoal.Position()
		angle := (180.0/math.Pi)*math.Atan2(y-body.Y, x-body.X) - body.T
		cmd += p.Client.Kick(60, angle)
	}
	return cmd
}
