package player

// State returns the state vector
func (p *Player) State() int {
	self := p.GetSelfData()
	state := (int(self.T) + 180)
	if state >= 360 {
		state -= 360
	}
	return state
}
