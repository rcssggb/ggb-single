package player

// Ball contains all the relevant information about the ball
type Ball struct {
	// Number of simulation cycles since ball was last seen
	NotSeenFor uint16

	// Relative polar coordinates
	Distance   float64
	Direction  float64
	DistChange float64
	DirChange  float64

	// Absolute cartesian coordinates
	X    float64
	Y    float64
	VelX float64
	VelY float64
}
