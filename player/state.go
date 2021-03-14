package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/rcsscommon"
)

// State returns the state vector
func (p *Player) State() int {
	self := p.GetSelfData()
	ball := p.GetBall()

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
	state := ballDistState
	shift := 7

	seesBall := 0
	if ball.NotSeenFor == 0 {
		seesBall = 1
	}
	state += seesBall * shift
	shift *= 2

	ballDir := ball.Direction
	const halfSizeBallDir = 30.0
	const nStatesBallDir = 5
	const stateSizeBallDir = 2 * halfSizeBallDir / nStatesBallDir
	if ballDir >= halfSizeBallDir {
		ballDir = halfSizeBallDir - 0.01
	} else if ballDir < -halfSizeBallDir {
		ballDir = -halfSizeBallDir
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
	const halfSizeX = rcsscommon.FieldMaxX + 5
	const nStatesX = 10
	const stateSizeX = 2 * halfSizeX / nStatesX
	if playerX > halfSizeX {
		playerX = halfSizeX - 0.01
	} else if playerX < -halfSizeX {
		playerX = -halfSizeX
	}
	playerX += halfSizeX
	playerXState := int(playerX / stateSizeX)
	if playerXState >= nStatesX {
		playerXState = nStatesX - 1
	} else if playerXState < 0 {
		playerXState = 0
	}
	state += playerXState * shift
	shift *= nStatesX

	playerY := self.Y
	const halfSizeY = rcsscommon.FieldMaxY + 5
	const nStatesY = 7
	const stateSizeY = 2 * halfSizeY / nStatesY
	if playerY > halfSizeY {
		playerY = halfSizeY - 0.01
	} else if playerY < -halfSizeY {
		playerY = -halfSizeY
	}
	playerY += halfSizeY
	playerYState := int(playerY / stateSizeY)
	if playerYState >= nStatesY {
		playerYState = nStatesY - 1
	} else if playerYState < 0 {
		playerYState = 0
	}
	state += playerYState * shift
	shift *= nStatesY

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
