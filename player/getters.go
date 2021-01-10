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

// GetSeenFriendly returns position of seen friendly players
func (p *Player) GetSeenFriendly() map[int]SeenPlayerPosition {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.friendlyPlayersPos
}

// GetSeenOpponent returns position of seen opponent players
func (p *Player) GetSeenOpponent() map[int]SeenPlayerPosition {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.opponentPlayersPos
}
