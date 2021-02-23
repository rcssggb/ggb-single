package player

// DiscreteAction takes one os 16 pre-defined discrete actions
func (p *Player) DiscreteAction(a int) string {
	switch a {
	case 0:
		return p.Client.Turn(0)
	case 1:
		return p.Client.Turn(10)
	case 2:
		return p.Client.Turn(-10)
	}
	return "(error invalid action index)"
}
