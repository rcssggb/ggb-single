package player

import (
	"math"
	"math/rand"

	"github.com/rcssggb/ggb-lib/rcsscommon"
)

// RandomPolicy performs a random action
func (p *Player) RandomPolicy() string {
	var cmd string
	action := rand.Intn(4)
	powerParam := rand.Float64()*200 - 100
	angleParam := rand.Float64()*360 - 180
	turnNeckAngle := rand.Float64()*360 - 180

	switch action {
	case 0:
	case 1:
		cmd += p.Client.Dash(powerParam, angleParam)
	case 2:
		cmd += p.Client.Kick(powerParam, angleParam)
	case 3:
		cmd += p.Client.Turn(angleParam)
	}
	cmd += p.Client.TurnNeck(turnNeckAngle)
	return cmd
}

// NaivePolicy performs the naive behavior (run to ball and kick to goal)
func (p *Player) NaivePolicy() string {
	var cmd string

	ball := p.GetBall()
	body := p.GetSelfData()

	if ball.NotSeenFor != 0 {
		cmd += p.Client.Turn(20)
	} else {
		ballDist := ball.Distance
		if ballDist < 0.7 {
			goalX, goalY := rcsscommon.FlagRightGoal.Position()
			goalAngle := (180.0/math.Pi)*math.Atan2(goalY-body.Y, goalX-body.X) - body.T
			cmd += p.Client.Kick(50, goalAngle)
		} else {
			ballAngle := ball.Direction + body.NeckAngle
			cmd += p.Client.Dash(30, ballAngle)
			cmd += p.Client.TurnNeck(ball.Direction / 2)
		}
	}

	return cmd
}
