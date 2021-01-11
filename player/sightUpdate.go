package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/rcsscommon"
)

// sightUpdate defines the goroutine that receives and
// processes visual information received by client
func (p *Player) sightUpdate() {
	lastTime := -1
	for {
		p.Client.WaitSight()
		p.mutex.Lock()
		data := p.Client.See()

		var selfAngleSamples []float64
		if data.Lines.Len() > 0 {
			var selfAngle float64
			closestLine := data.Lines[0]
			lDir := closestLine.Direction
			if lDir < 0 {
				lDir += 90
			} else {
				lDir -= 90
			}
			switch closestLine.ID {
			case rcsscommon.LineRight:
				selfAngle = 0 - lDir
			case rcsscommon.LineBottom:
				selfAngle = 90 - lDir
			case rcsscommon.LineLeft:
				selfAngle = 180 - lDir
			case rcsscommon.LineTop:
				selfAngle = -90 - lDir
			}

			// If you see 2 lines it means you're outside the field
			if data.Lines.Len() >= 2 {
				selfAngle += 180
			}

			selfAngle -= p.self.NeckAngle

			if selfAngle > 180 {
				selfAngle -= 360
			} else if selfAngle < -180 {
				selfAngle += 360
			}

			selfAngleSamples = append(selfAngleSamples, selfAngle)
		}
		// TODO: improve number of samples selfAngleSamples when seeing 2+ flags

		if len(selfAngleSamples) > 0 {
			// TODO: calculate sample average
			p.self.T = selfAngleSamples[0]
		}

		if data.Flags.Len() > 0 {
			var xAcc, yAcc float64 = 0, 0
			for _, f := range data.Flags {
				xFlag, yFlag := f.ID.Position()
				absAngle := (math.Pi / 180.0) * (f.Direction + p.self.T + p.self.NeckAngle)
				xTmp := xFlag - math.Cos(absAngle)*f.Distance
				yTmp := yFlag - math.Sin(absAngle)*f.Distance
				xAcc += xTmp
				yAcc += yTmp
			}
			p.self.X = xAcc / (float64)(data.Flags.Len())
			p.self.Y = yAcc / (float64)(data.Flags.Len())
		}

		if data.Ball != nil {
			ball := *data.Ball
			p.ball.NotSeenFor = 0

			// Relative coordinates
			p.ball.Distance = ball.Distance
			p.ball.Direction = ball.Direction
			p.ball.DistChange = ball.DistChange
			p.ball.DirChange = ball.DirChange

			/* Absolute coordinates */
			// Calculate sin and cos of vector from player to object
			sin, cos := math.Sincos(math.Pi / 180.0 * (p.self.T + p.self.NeckAngle - p.ball.Direction))
			// Project to absolute frame of reference
			p.ball.X = p.self.X + ball.Distance*cos
			p.ball.Y = p.self.Y + ball.Distance*sin
			// Multiply DirChange by relative vector length and
			// rotate the vectors to the absolute frame of reference
			p.ball.VelX = ball.DistChange*cos - ball.DirChange*ball.Distance*sin + p.self.VelX
			p.ball.VelY = ball.DistChange*sin + ball.DirChange*ball.Distance*cos + p.self.VelY
		} else {
			// If ball was not seen, increment NotSeenFor timer
			if data.Time > lastTime {
				p.ball.NotSeenFor += (uint16)(data.Time - lastTime)
			}
		}

		// Before updating players according to what was seen,
		// increment NotSeenFor timer for every player
		// (seen players will have timer reset later)
		if data.Time > lastTime {
			for u := range p.friendlyPlayersPos {
				pPos := p.friendlyPlayersPos[u]
				pPos.NotSeenFor += (uint16)(data.Time - lastTime)
				p.friendlyPlayersPos[u] = pPos
			}
			for u := range p.opponentPlayersPos {
				pPos := p.opponentPlayersPos[u]
				pPos.NotSeenFor += (uint16)(data.Time - lastTime)
				p.opponentPlayersPos[u] = pPos
			}
		}

		// Update seen players positions
		for _, seenPlayer := range data.Players {
			if seenPlayer.Unum == -1 {
				// For now we are ignoring not fully known players
				break
			}
			var seenPlayerPos = SeenPlayerPosition{
				Distance:   seenPlayer.Distance,
				Direction:  seenPlayer.Direction,
				DistChange: seenPlayer.DistChange,
				DirChange:  seenPlayer.DirChange,
				BodyDir:    seenPlayer.BodyDir + p.self.T + p.self.NeckAngle,
				NeckDir:    seenPlayer.NeckDir + p.self.T + p.self.NeckAngle,
			}

			// Calculate sin and cos of vector from player to object
			sin, cos := math.Sincos(math.Pi / 180.0 * (seenPlayer.Direction - (p.self.T + p.self.NeckAngle)))
			// Project to absolute frame of reference
			seenPlayerPos.X = p.self.X + seenPlayer.Distance*cos
			seenPlayerPos.Y = p.self.Y + seenPlayer.Distance*sin
			// Multiply DirChange by relative vector length and
			// rotate the vectors to the absolute frame of reference
			seenPlayerPos.VelX = seenPlayer.DistChange*cos - seenPlayer.DirChange*seenPlayer.Distance*sin + p.self.VelX
			seenPlayerPos.VelY = seenPlayer.DistChange*sin + seenPlayer.DirChange*seenPlayer.Distance*cos + p.self.VelY

			if seenPlayer.Team == p.Client.TeamName() {
				p.friendlyPlayersPos[seenPlayer.Unum] = seenPlayerPos
			} else {
				p.opponentPlayersPos[seenPlayer.Unum] = seenPlayerPos
			}
		}

		if data.Time > lastTime {
			lastTime = data.Time
		}

		p.mutex.Unlock()
	}
}
