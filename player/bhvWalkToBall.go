package player

func (p *Player) bhvWalkToBall() string {
	cmd := ""
	ball := p.GetBall()
	body := p.GetSelfData()
	ballAngle := ball.Direction + body.NeckAngle
	if -15 < ballAngle && ballAngle < 15 {
		cmd += p.Client.Dash(60, ballAngle)
	} else {
		cmd += p.Client.Turn(ballAngle)
	}
	return cmd
}
