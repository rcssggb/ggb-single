package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/rcsscommon"
)

func (p *Player) bhvLocateBall() string {
	cmd := ""
	ball := p.GetBall()
	self := p.GetSelfData()
	if ball.NotSeenFor > 0 {
		lastBallAngle := math.Atan2(ball.Y-self.Y, ball.X-self.X) * (180.0 / math.Pi)
		amountToTurn := lastBallAngle - self.T + self.NeckAngle
		if ball.NotSeenFor > 30 {
			amountToTurn = 45
		} else {
			if amountToTurn > 180 {
				amountToTurn -= 360
			} else if amountToTurn < -180 {
				amountToTurn += 360
			}
		}
		cmd += p.Client.Turn(amountToTurn)
	} else {
		ballAngle := ball.Direction
		amountToTurn := ballAngle / 5
		amountToTurnNeck := ballAngle / 2
		cmd += p.Client.Turn(amountToTurn)
		cmd += p.Client.TurnNeck(amountToTurnNeck)
	}
	return cmd
}

func (p *Player) bhvLeadBallRight() string {
	cmd := ""
	body := p.GetSelfData()
	ball := p.GetBall()
	if ball.Distance < 0.8 {
		x := body.X + 10
		y := body.Y
		angle := (180.0/math.Pi)*math.Atan2(y-body.Y, x-body.X) - body.T
		cmd += p.Client.Kick(5, angle)
	}
	return cmd
}

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

func (p *Player) bhvShootToGoalL() string {
	cmd := ""
	body := p.GetSelfData()
	ball := p.GetBall()
	if ball.Distance < 0.8 {
		x, y := rcsscommon.FlagLeftGoal.Position()
		angle := (180.0/math.Pi)*math.Atan2(y-body.Y, x-body.X) - body.T
		cmd += p.Client.Kick(30, angle)
	}
	return cmd
}

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

func (p *Player) bhvWalkAwayFromBall() string {
	cmd := ""
	ball := p.GetBall()
	body := p.GetSelfData()
	if ball.NotSeenFor == 0 {
		ballAngle := ball.Direction + body.NeckAngle
		cmd += p.Client.Dash(-60, ballAngle)
	}
	return cmd
}

func (p *Player) bhvWalkToBall() string {
	cmd := ""
	ball := p.GetBall()
	body := p.GetSelfData()
	if ball.NotSeenFor == 0 {
		ballAngle := ball.Direction + body.NeckAngle
		cmd += p.Client.Dash(60, ballAngle)
	}
	return cmd
}
