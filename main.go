package main

import (
	"log"
	"time"

	"github.com/rcssggb/ggb-lib/trainerclient"
	"github.com/rcssggb/ggb-single/player"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	hostName := "rcssserver"

	for {
		p, err := player.NewPlayer("single-agent", hostName)
		if err != nil {
			log.Println(err)
			continue
		}

		t, err := trainerclient.NewTrainerClient(hostName)
		if err != nil {
			log.Println(err)
			continue
		}

		t.EarOn()
		t.EyeOn()

		time.Sleep(2 * time.Second)

		t.Start()
		for {
			currentTime := p.Client.Time()

			if currentTime != 0 {
				// time.Sleep(10 * time.Millisecond)
			}

			p.NaivePolicy()

			if p.Client.PlayMode().String() == "time_over" {
				p.Client.Bye()
				break
			}

			err := p.Client.Error()
			for err != nil {
				p.Client.Log(err)
				err = p.Client.Error()
			}

			t.DoneSynch()
			p.WaitCycle()
		}

		time.Sleep(5 * time.Second)
	}
}
