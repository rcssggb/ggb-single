package player

// DiscreteAction takes one os 16 pre-defined discrete actions
func (p *Player) DiscreteAction(a int) string {
	switch a {
	case 0:
		return p.Client.Turn(0)
	case 1:
		return p.Client.Turn(7)
	case 2:
		return p.Client.Turn(-7)
	case 3:
		return p.Client.Turn(15)
	case 4:
		return p.Client.Turn(-15)
	case 5:
		return p.Client.Dash(30, 0)
	case 6:
		return p.Client.Kick(25, 0)
	}
	return "(error invalid action index)"
}
