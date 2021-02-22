package player

import (
	"math"
)

func (p *Player) bhvLocateBall() string {
	cmd := ""
	ball := p.GetBall()
	self := p.GetSelfData()
	if ball.NotSeenFor > 0 {
		lastBallAngle := math.Atan2(ball.Y-self.Y, ball.X-self.X) * (180.0 / math.Pi)
		amountToTurn := lastBallAngle - self.T + self.NeckAngle
		if amountToTurn > 180 {
			amountToTurn -= 360
		} else if amountToTurn < -180 {
			amountToTurn += 360
		}
		cmd += p.Client.Turn(amountToTurn)
	} else if ball.NotSeenFor == 0 {
		ballAngle := ball.Direction
		amountToTurn := ballAngle / 5
		amountToTurnNeck := ballAngle / 2
		cmd += p.Client.Turn(amountToTurn)
		cmd += p.Client.TurnNeck(amountToTurnNeck)
	} else {
		cmd += p.Client.Turn(30)
	}
	return cmd
}
