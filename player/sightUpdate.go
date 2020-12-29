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
				p.selfPos.T = 0 - lDir
			case rcsscommon.LineBottom:
				p.selfPos.T = 90 - lDir
			case rcsscommon.LineLeft:
				p.selfPos.T = 180 - lDir
			case rcsscommon.LineTop:
				p.selfPos.T = -90 - lDir
			}
		}

		// If you see 2 lines it means you're outside the field
		if data.Lines.Len() >= 2 {
			p.selfPos.T += 180
		}

		p.selfPos.T -= p.body.NeckAngle

		if p.selfPos.T > 180 {
			p.selfPos.T -= 360
		} else if p.selfPos.T < -180 {
			p.selfPos.T += 360
		}

		if data.Flags.Len() > 0 {
			var xAcc, yAcc float64 = 0, 0
			for _, f := range data.Flags {
				xFlag, yFlag := f.ID.Position()
				absAngle := (3.14159 / 180.0) * (f.Direction + p.selfPos.T + p.body.NeckAngle)
				xTmp := xFlag - math.Cos(absAngle)*f.Distance
				yTmp := yFlag - math.Sin(absAngle)*f.Distance
				xAcc += xTmp
				yAcc += yTmp
			}
			p.selfPos.X = xAcc / (float64)(data.Flags.Len())
			p.selfPos.Y = yAcc / (float64)(data.Flags.Len())
		}

		if data.Ball != nil {
			// Relative coordinates
			p.ball.Distance = data.Ball.Distance
			p.ball.Direction = data.Ball.Direction - p.body.NeckAngle
			p.ball.DistChange = data.Ball.DistChange
			p.ball.DirChange = data.Ball.DirChange

			// Absolute coordinates
			sin, cos := math.Sincos(math.Pi / 180.0 * (p.ball.Direction - p.selfPos.T))
			p.ball.X = p.selfPos.X + p.ball.Distance*cos
			p.ball.Y = p.selfPos.Y + p.ball.Distance*sin
			p.ball.VelX = p.ball.DistChange*cos + p.ball.DirChange*p.ball.Distance*sin
			p.ball.VelY = p.ball.DistChange*sin + p.ball.DirChange*p.ball.Distance*cos

			p.ball.NotSeenFor = 0
		} else {
			// TODO: update ball old relative coordinates considering new selfpos?
			if data.Time > lastTime {
				p.ball.NotSeenFor++
			}
		}

		if len(data.Players) > 0 {
			for _, player := range data.Players {
				pl := PlayerPosition{
					Team:        player.Team,
					Unum:        player.Unum,
					Distance:    player.Distance,
					Direction:   player.Direction,
					DistChange:  player.DistChange,
					DirChange:   player.DirChange,
					BodyDir:     player.BodyDir,
					NeckDir:     player.NeckDir,
					IsPointing:  player.IsPointing,
					PointingDir: player.PointingDir,
					Action:      player.Action,
				}
				p.PlayersPos = append(p.PlayersPos, pl)
			}
		}

		if data.Time > lastTime {
			lastTime = data.Time
		}
	}
}
