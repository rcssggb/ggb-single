package player

// WaitCycle waits for next simulation cycle
func (p *Player) WaitCycle() {
	if p.Client.ServerParams().SynchMode {
		p.Client.DoneSynch()
		p.Client.WaitSynch()
	} else {
		p.Client.WaitNextStep(p.Client.Time())
	}
}
