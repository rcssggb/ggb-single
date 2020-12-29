package player

// PlayerPosition contains data about seen players.
// Team == "" means you don't know the player's team.
// Unum == -1 means you don't know the player's shirt num.
type PlayerPosition struct {
	Team string
	Unum int

	Distance   float64
	Direction  float64
	DistChange float64
	DirChange  float64
	BodyDir    float64
	NeckDir    float64

	IsPointing  bool
	PointingDir float64
	Action      string
}
