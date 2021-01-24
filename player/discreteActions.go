package player

// DiscreteAction takes one os 16 pre-defined discrete actions
func (p *Player) DiscreteAction(a int) string {
	switch a {
	case 0:
		return p.Client.Dash(50, 0)
	case 1:
		return p.Client.Dash(20, 0)
	case 2:
		return p.Client.Dash(80, 0)
	case 3:
		return p.Client.Dash(-50, 0)
	case 4:
		return p.Client.Kick(50, 0)
	case 5:
		return p.Client.Kick(50, 45)
	case 6:
		return p.Client.Kick(50, -45)
	case 7:
		return p.Client.Kick(-50, 0)
	case 8:
		return p.Client.Turn(30)
	case 9:
		return p.Client.Turn(-30)
	case 10:
		return p.Client.Turn(90)
	case 11:
		return p.Client.Turn(-90)
	case 12:
		return p.Client.Kick(20, 0)
	case 13:
		return p.Client.Kick(20, 45)
	case 14:
		return p.Client.Kick(20, -45)
	case 15:
		return p.Client.Move(-30, 0)
	}
	return "(error invalid action index)"
}
