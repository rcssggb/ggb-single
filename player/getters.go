package player

// GetSelfData returns the current player position
func (p Player) GetSelfData() SelfData {
	return p.self
}

// GetBall returns the current ball info
func (p Player) GetBall() Ball {
	return p.ball
}
