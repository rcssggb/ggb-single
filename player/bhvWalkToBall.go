package player

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
