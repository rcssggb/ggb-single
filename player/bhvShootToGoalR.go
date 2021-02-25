package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/rcsscommon"
)

func (p *Player) bhvShootToGoalR() string {
	cmd := ""
	body := p.GetSelfData()
	x, y := rcsscommon.FlagRightGoal.Position()
	angle := (180.0/math.Pi)*math.Atan2(y-body.Y, x-body.X) - body.T
	cmd += p.Client.Kick(30, angle)
	return cmd
}
