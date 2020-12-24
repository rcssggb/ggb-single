package player

import (
	"time"
)

// bodyUpdate defines the goroutine that receives and
// processes body sensor information received by client
func (p *Player) bodyUpdate() {
	currentTimestamp := -1
	for {
		data := p.Client.SenseBody()

		// This if statement and whole timing logic must be changed after ggb-lib implements new data signaling
		if data.Time <= currentTimestamp {
			time.Sleep(10 * time.Millisecond)
			if data.Time != currentTimestamp {
				continue
			}
		}

		currentTimestamp = data.Time
		p.body.NeckAngle = data.HeadAngle
	}
}
