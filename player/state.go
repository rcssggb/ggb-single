package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/rcsscommon"
)

// State returns the state vector
func (p *Player) State() int {
	self := p.GetSelfData()
	ball := p.GetBall()

	seesBall := false
	if ball.NotSeenFor == 0 {
		seesBall = true
	}

	var ballDist float64
	if seesBall {
		ballDist = ball.Distance
	} else {
		ballDist = math.Sqrt(math.Pow(ball.X-self.X, 2) + math.Pow(ball.Y-self.Y, 2))
	}
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
	state := ballDistState
	shift := 7

	const halfSizeBallDir = 180.0
	const nStatesBallDir = 24
	const stateSizeBallDir = 2 * halfSizeBallDir / nStatesBallDir
	var ballDir float64
	if seesBall {
		ballDir = ball.Direction
	} else {
		ballDir = 180.0/math.Pi*math.Atan2(ball.Y-self.Y, ball.X-self.X) - self.T
	}
	if ballDir >= halfSizeBallDir {
		ballDir -= halfSizeBallDir
	} else if ballDir < -halfSizeBallDir {
		ballDir += halfSizeBallDir
	}
	ballDir += halfSizeBallDir
	ballDirState := int(ballDir / stateSizeBallDir)
	if ballDirState >= nStatesBallDir {
		ballDirState = nStatesBallDir - 1
	} else if ballDirState < 0 {
		ballDirState = 0
	}
	state += ballDirState * shift
	shift *= nStatesBallDir

	playerX := self.X
	const halfSizePlayerX = rcsscommon.FieldMaxX + 5
	const nStatesPlayerX = 10
	const stateSizePlayerX = 2 * halfSizePlayerX / nStatesPlayerX
	if playerX > halfSizePlayerX {
		playerX = halfSizePlayerX - 0.01
	} else if playerX < -halfSizePlayerX {
		playerX = -halfSizePlayerX
	}
	playerX += halfSizePlayerX
	playerXState := int(playerX / stateSizePlayerX)
	if playerXState >= nStatesPlayerX {
		playerXState = nStatesPlayerX - 1
	} else if playerXState < 0 {
		playerXState = 0
	}
	state += playerXState * shift
	shift *= nStatesPlayerX

	playerY := self.Y
	const halfSizePlayerY = rcsscommon.FieldMaxY + 5
	const nStatesPlayerY = 7
	const stateSizePlayerY = 2 * halfSizePlayerY / nStatesPlayerY
	if playerY > halfSizePlayerY {
		playerY = halfSizePlayerY - 0.01
	} else if playerY < -halfSizePlayerY {
		playerY = -halfSizePlayerY
	}
	playerY += halfSizePlayerY
	playerYState := int(playerY / stateSizePlayerY)
	if playerYState >= nStatesPlayerY {
		playerYState = nStatesPlayerY - 1
	} else if playerYState < 0 {
		playerYState = 0
	}
	state += playerYState * shift
	shift *= nStatesPlayerY

	playerT := self.T
	const halfSizeT = 180.0
	const nStatesT = 24
	const stateSizeT = 2 * halfSizeT / nStatesT
	if playerT > halfSizeT {
		playerT = halfSizeT - 0.01
	} else if playerT < -halfSizeT {
		playerT = -halfSizeT
	}
	playerT += halfSizeT
	playerTState := int(playerT / stateSizeT)
	if playerTState >= nStatesT {
		playerTState = nStatesT - 1
	} else if playerTState < 0 {
		playerTState = 0
	}
	state += playerTState * shift
	shift *= nStatesT

	return state
}
