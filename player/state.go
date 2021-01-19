package player

// State returns the state vector
func (p *Player) State() []float64 {
	self := p.GetSelfData()
	ball := p.GetBall()
	return []float64{
		self.Stamina,
		self.Effort,
		self.Capacity,
		self.X,
		self.Y,
		self.T,
		self.VelX,
		self.VelY,
		self.NeckAngle,
		self.VelSpeed,
		self.VelDir,
		float64(ball.NotSeenFor),
		ball.Distance,
		ball.Direction,
		ball.DistChange,
		ball.DirChange,
		ball.X,
		ball.Y,
		ball.VelX,
		ball.VelY,
		// TODO: encode play_mode
	}
}
