package player

import "math"

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

	ballDist := ball.Distance
	if ballDist < 1 {
		ballDist = 1
	}
	ballDistState := int(math.Log2(ballDist))
	if ballDistState > 5 {
		ballDistState = 5
	} else if ballDistState < 0 {
		ballDistState = 0
	}
	state += ballDistState * shift
	shift *= 6

	seesBall := 0
	if ball.NotSeenFor == 0 {
		seesBall = 1
	}
	state += seesBall * shift
	shift *= 2

	return state
}
