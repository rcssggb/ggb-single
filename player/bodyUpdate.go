package player

// bodyUpdate defines the goroutine that receives and
// processes body sensor information received by client
func (p *Player) bodyUpdate() {
	for {
		p.Client.WaitBody()

		data := p.Client.SenseBody()

		p.body.NeckAngle = data.HeadAngle
	}
}
