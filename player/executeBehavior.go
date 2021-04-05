package player

// ExecuteBehavior executes b from pre-defined behaviors
func (p *Player) ExecuteBehavior(b int) string {
	// fmt.Println("behavior: ", b)
	switch b {
	case 0:
		return ""
	case 1:
		return p.bhvLocateBall()
	case 2:
		return p.bhvLeadBallRight()
	case 3:
		return p.bhvShootToGoalR()
	case 4:
		return p.bhvWalkToBall()
	case 5:
		return p.bhvLeadBallLeft()
	case 6:
		return p.bhvShootToGoalL()
	}
	return "(error invalid behavior index)"
}
