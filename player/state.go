package player

// State returns the state vector
func (p *Player) State() int {
	// self := p.GetSelfData()
	ball := p.GetBall()

	ballDirState := (int(ball.Direction) + 30) / 5
	if ballDirState >= 12 {
		ballDirState = 11
	} else if ballDirState < 0 {
		ballDirState = 0
	}
	state := ballDirState
	shift := 12

	seesBall := 0
	if ball.NotSeenFor == 0 {
		seesBall = 1
	}
	state += seesBall * shift
	shift *= 2

	return state
}
