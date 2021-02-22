package player

// ExecuteBehavior executes b from pre-defined behaviors
func (p *Player) ExecuteBehavior(b int) string {
	switch b {
	case 0:
		return p.bhvLocateBall()
	case 1:
		return p.bhvLeadBall()
	}
	return "(error invalid behavior index)"
}
