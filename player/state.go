package player

// State returns the state vector
func (p *Player) State() int {
	self := p.GetSelfData()
	state := (int(self.T) + 180) / 15
	if state > 23 {
		state = 23
	} else if state < 0 {
		state = 0
	}
	return state
}
