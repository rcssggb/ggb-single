package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/rcsscommon"
)

// State returns the state vector
func (p *Player) State() int {
	self := p.GetSelfData()
	ball := p.GetBall()

	ballDir := ball.Direction
	if ballDir >= 30 {
		ballDir = 29.99
	} else if ballDir < -30 {
		ballDir = -30
	}
	ballDir += 30
	ballDirState := int(ballDir) / 5
	if ballDirState > 4 {
		ballDirState = 4
	} else if ballDirState < 0 {
		ballDirState = 0
	}
	state := ballDirState
	shift := 5

	ballDist := ball.Distance
	ballDistScaleFactor := 0.7
	ballDist /= ballDistScaleFactor
	if ballDist < 1 {
		ballDist = 1
	}
	ballDistState := int(math.Log2(ballDist))
	if ballDistState > 6 {
		ballDistState = 6
	} else if ballDistState < 0 {
		ballDistState = 0
	}
	state += ballDistState * shift
	shift *= 7

	playerX := self.X
	if playerX > rcsscommon.FieldMaxX {
		playerX = rcsscommon.FieldMaxX - 0.01
	} else if playerX < -rcsscommon.FieldMaxX {
		playerX = -rcsscommon.FieldMaxX
	}
	playerX += rcsscommon.FieldMaxX
	playerXState := int(playerX) / 10
	if playerXState > 9 {
		playerXState = 9
	} else if playerXState < 0 {
		playerXState = 0
	}
	state += playerXState * shift
	shift *= 10

	// playerY := self.Y
	// if playerY > rcsscommon.FieldMaxY {
	// 	playerY = rcsscommon.FieldMaxY - 0.01
	// } else if playerY < -rcsscommon.FieldMaxY {
	// 	playerY = -rcsscommon.FieldMaxY
	// }
	// playerY += rcsscommon.FieldMaxY
	// playerYState := int(playerY) / 7
	// if playerYState > 6 {
	// 	playerYState = 6
	// } else if playerYState < 0 {
	// 	playerYState = 0
	// }
	// state += playerYState * shift
	// shift *= 7

	playerT := int((self.T + 180)) / 12
	if playerT > 11 {
		playerT = 11
	} else if playerT < 0 {
		playerT = 0
	}
	state += playerT * shift
	shift *= 12

	seesBall := 0
	if ball.NotSeenFor == 0 {
		seesBall = 1
	}
	state += seesBall * shift
	shift *= 2

	return state
}
