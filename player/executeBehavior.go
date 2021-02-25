package player

// ExecuteBehavior executes b from pre-defined behaviors
func (p *Player) ExecuteBehavior(b int) string {
	// fmt.Println("behavior: ", b)
	switch b {
	case 0:
		return p.bhvLocateBall()
	case 1:
		return p.bhvLeadBallRight()
	case 2:
		return p.bhvShootToGoalR()
	case 3:
		return p.bhvWalkToBall()
	case 4:
		return p.bhvLeadBallLeft()
	case 5:
		return p.bhvShootToGoalL()
	}
	return "(error invalid behavior index)"
}
