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
		return p.bhvWalkToBall()
	case 3:
		return p.bhvLeadBallRight()
	case 4:
		return p.bhvShootToGoalR()
	case 5:
		return p.bhvWalkAwayFromBall()
	case 6:
		return p.bhvLeadBallLeft()
	case 7:
		return p.bhvShootToGoalL()
	}
	return "(error invalid behavior index)"
}
