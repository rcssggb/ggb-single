package player

import "math"

// bodyUpdate defines the goroutine that receives and
// processes body sensor information received by client
func (p *Player) bodyUpdate() {
	for {
		p.Client.WaitBody()
		p.mutex.Lock()

		data := p.Client.SenseBody()

		p.self.NeckAngle = data.HeadAngle
		p.self.VelDir = data.Speed.Direction
		p.self.VelSpeed = data.Speed.Magnitude

		sin, cos := math.Sincos(math.Pi / 180.0 * (p.self.T - p.self.VelDir))
		p.self.VelX = p.self.VelSpeed * cos
		p.self.VelY = p.self.VelSpeed * sin

		p.mutex.Unlock()
	}
}
