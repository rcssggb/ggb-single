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
