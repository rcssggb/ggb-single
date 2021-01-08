package player

// PlayerPosition contains data about seen players.
// Team == "" means you don't know the player's team.
type PlayerPosition struct {
	// Number of simulation cycles since player was last seen
	NotSeenFor uint16

	// Relative polar coordinates
	Distance   float64
	Direction  float64
	DistChange float64
	DirChange  float64

	// Absolute polar coordinates
	BodyDir float64
	NeckDir float64

	// Absolute cartesian coordinates
	X    float64
	Y    float64
	VelX float64
	VelY float64
}
