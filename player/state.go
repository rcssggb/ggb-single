package player

// State returns the state vector
func (p *Player) State() []float64 {
	// self := p.GetSelfData()
	ball := p.GetBall()
	serverParams := p.Client.ServerParams()
	// playModeOneHot := p.Client.PlayMode().OneHot()
	ret := []float64{
		// self.Stamina / serverParams.StaminaMax,
		// self.Effort,
		// self.Capacity / serverParams.StaminaCapacity,
		// self.X / rcsscommon.FieldMaxX,
		// self.Y / rcsscommon.FieldMaxY,
		// self.T / 180.0,
		// self.VelX / serverParams.PlayerSpeedMax,
		// self.VelY / serverParams.PlayerSpeedMax,
		// self.NeckAngle / serverParams.MaxNeckAng,
		// self.VelSpeed / serverParams.PlayerSpeedMax,
		// self.VelDir / 180.0,
		// float64(ball.NotSeenFor),
		ball.Distance / serverParams.VisibleDistance,
		ball.Direction / 180.0,
		// ball.DistChange / serverParams.BallSpeedMax,
		// ball.DirChange, // TODO: don't know how to normalize this one
		// ball.X / rcsscommon.FieldMaxX,
		// ball.Y / rcsscommon.FieldMaxY,
		// ball.VelX / serverParams.BallSpeedMax,
		// ball.VelY / serverParams.BallSpeedMax,
	}
	return ret
}
