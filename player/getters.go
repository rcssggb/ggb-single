package player

// GetSelfPos returns the current player position
func (p Player) GetSelfPos() Position {
	return p.selfPos
}

// GetBody returns the current player body info
func (p Player) GetBody() Body {
	return p.body
}
