package player

// DiscreteAction takes one os 16 pre-defined discrete actions
func (p *Player) DiscreteAction(a int) string {
	switch a {
	case 0:
		return ""
	case 1:
		return p.Client.Turn(7)
	case 2:
		return p.Client.Turn(-7)
	case 3:
		return p.Client.Turn(15)
	case 4:
		return p.Client.Turn(-15)
	case 5:
		return p.Client.Turn(31)
	case 6:
		return p.Client.Turn(-31)
	case 7:
		return p.Client.Dash(50, 0)
	case 8:
		return p.Client.Dash(50, 30)
	case 9:
		return p.Client.Dash(50, -30)
	case 10:
		return p.Client.Kick(50, 0)
	case 11:
		return p.Client.Kick(50, 45)
	case 12:
		return p.Client.Kick(50, -45)
	}
	return "(error invalid action index)"
}
