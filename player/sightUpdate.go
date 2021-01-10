package player

import (
	"math"

	"github.com/rcssggb/ggb-lib/playerclient/parser"
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

		var closestLine *parser.LineData
		closestLine = nil
		if data.Lines.Len() > 0 {
			closestLine = &data.Lines[0]
			lDir := closestLine.Direction
			if lDir < 0 {
				lDir += 90
			} else {
				lDir -= 90
			}
			switch closestLine.ID {
			case rcsscommon.LineRight:
				p.self.T = 0 - lDir
			case rcsscommon.LineBottom:
				p.self.T = 90 - lDir
			case rcsscommon.LineLeft:
				p.self.T = 180 - lDir
			case rcsscommon.LineTop:
				p.self.T = -90 - lDir
			}
		}

		// If you see 2 lines it means you're outside the field
		if data.Lines.Len() >= 2 {
			p.self.T += 180
		}

		p.self.T -= p.self.NeckAngle

		if p.self.T > 180 {
			p.self.T -= 360
		} else if p.self.T < -180 {
			p.self.T += 360
		}

		if data.Flags.Len() > 0 {
			var xAcc, yAcc float64 = 0, 0
			for _, f := range data.Flags {
				xFlag, yFlag := f.ID.Position()
				absAngle := (3.14159 / 180.0) * (f.Direction + p.self.T + p.self.NeckAngle)
				xTmp := xFlag - math.Cos(absAngle)*f.Distance
				yTmp := yFlag - math.Sin(absAngle)*f.Distance
				xAcc += xTmp
				yAcc += yTmp
			}
			p.self.X = xAcc / (float64)(data.Flags.Len())
			p.self.Y = yAcc / (float64)(data.Flags.Len())
		}

		if data.Ball != nil {
			p.ball.NotSeenFor = 0

			// Relative coordinates
			p.ball.Distance = data.Ball.Distance
			p.ball.Direction = data.Ball.Direction - p.self.NeckAngle
			p.ball.DistChange = data.Ball.DistChange
			p.ball.DirChange = data.Ball.DirChange

			/* Absolute coordinates */
			// Calculate sin and cos of vector from player to object
			sin, cos := math.Sincos(math.Pi / 180.0 * (p.self.T - p.ball.Direction))
			// Project to absolute frame of reference
			p.ball.X = p.self.X + p.ball.Distance*cos
			p.ball.Y = p.self.Y + p.ball.Distance*sin
			// Multiply DirChange by relative vector length and
			// rotate the vectors to the absolute frame of reference
			p.ball.VelX = p.ball.DistChange*cos - p.ball.DirChange*p.ball.Distance*sin - p.self.VelX
			p.ball.VelY = p.ball.DistChange*sin + p.ball.DirChange*p.ball.Distance*cos - p.self.VelY
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
			sin, cos := math.Sincos(math.Pi / 180.0 * (seenPlayer.Direction - p.self.T))
			// Project to absolute frame of reference
			seenPlayerPos.X = p.self.X + seenPlayer.Distance*cos
			seenPlayerPos.Y = p.self.Y + seenPlayer.Distance*sin
			// Multiply DirChange by relative vector length and
			// rotate the vectors to the absolute frame of reference
			seenPlayerPos.VelX = seenPlayer.DistChange*cos - seenPlayer.DirChange*seenPlayer.Distance*sin
			seenPlayerPos.VelY = seenPlayer.DistChange*sin + seenPlayer.DirChange*seenPlayer.Distance*cos

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
