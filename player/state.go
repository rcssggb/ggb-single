package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/rcsscommon"
	"gonum.org/v1/gonum/mat"
)

// State returns the state vector
func (p *Player) State() mat.Vector {
	self := p.GetSelfData()
	ball := p.GetBall()

	seesBall := false
	if ball.NotSeenFor == 0 {
		seesBall = true
	}

	state := make([]float64, 0, 5)

	var ballDist float64
	if seesBall {
		ballDist = ball.Distance
	} else {
		ballDist = math.Sqrt(math.Pow(ball.X-self.X, 2) + math.Pow(ball.Y-self.Y, 2))
	}
	state = append(state, ballDist)

	const halfSizeBallDir = 180.0
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
	ballDir /= halfSizeBallDir
	state = append(state, ballDir)

	playerX := self.X
	const halfSizePlayerX = rcsscommon.FieldMaxX + 5
	if playerX > halfSizePlayerX {
		playerX = halfSizePlayerX
	} else if playerX < -halfSizePlayerX {
		playerX = -halfSizePlayerX
	}
	playerX /= halfSizePlayerX
	state = append(state, playerX)

	playerY := self.Y
	const halfSizePlayerY = rcsscommon.FieldMaxY + 5
	if playerY > halfSizePlayerY {
		playerY = halfSizePlayerY
	} else if playerY < -halfSizePlayerY {
		playerY = -halfSizePlayerY
	}
	playerY /= halfSizePlayerY
	state = append(state, playerY)

	playerT := self.T
	const halfSizeT = 180.0
	if playerT >= halfSizeT {
		playerT -= halfSizeT
	} else if playerT < -halfSizeT {
		playerT += halfSizeT
	}
	playerT /= halfSizeT
	state = append(state, playerT)

	stateVector := mat.NewVecDense(len(state), state)

	return stateVector
}
