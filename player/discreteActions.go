package player

import "math/rand"

// SampleDiscreteActionVector samples a random discrete action vector
func (p *Player) SampleDiscreteActionVector() []float64 {
	a := rand.Intn(16)
	vec := make([]float64, 16)
	vec[a] = 1
	return vec
}

// DiscreteActionVector generates an action vector
func (p *Player) DiscreteActionVector(a int) []float64 {
	vec := make([]float64, 16)
	vec[a] = 1
	return vec
}

// DiscreteAction takes one os 16 pre-defined discrete actions
func (p *Player) DiscreteAction(a int) {
	switch a {
	case 0:
		p.Client.Dash(50, 0)
	case 1:
		p.Client.Dash(20, 0)
	case 2:
		p.Client.Dash(80, 0)
	case 3:
		p.Client.Dash(-50, 0)
	case 4:
		p.Client.Kick(50, 0)
	case 5:
		p.Client.Kick(50, 45)
	case 6:
		p.Client.Kick(50, -45)
	case 7:
		p.Client.Kick(-50, 0)
	case 8:
		p.Client.Turn(30)
	case 9:
		p.Client.Turn(-30)
	case 10:
		p.Client.Turn(90)
	case 11:
		p.Client.Turn(-90)
	case 12:
		p.Client.Kick(20, 0)
	case 13:
		p.Client.Kick(20, 45)
	case 14:
		p.Client.Kick(20, -45)
	case 15:
		// do nothing
	}
}
