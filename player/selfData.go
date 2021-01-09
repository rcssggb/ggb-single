package player

// SelfData is a struct containing all relevant data about the player itself
type SelfData struct {
	// Absolute coordinates
	X float64
	Y float64
	T float64

	// Relative Coordinates
	NeckAngle float64
	VelSpeed  float64
	VelDir    float64
}
