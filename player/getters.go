package player

// GetSelfData returns the current player position
func (p *Player) GetSelfData() SelfData {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.self
}

// GetBall returns the current ball info
func (p *Player) GetBall() Ball {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.ball
}
