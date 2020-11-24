package main

import (
	"log"
	"time"

	"github.com/rcssggb/ggb-lib/playerclient"
	"github.com/rcssggb/ggb-lib/rcsscommon"
	"github.com/rcssggb/ggb-lib/trainerclient"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	hostName := "rcssserver"

	p, err := playerclient.NewPlayerClient("single-agent", hostName)
	if err != nil {
		log.Fatalln(err)
	}

	trainer, err := trainerclient.NewTrainerClient(hostName)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(2 * time.Second)

	trainer.EyeOn()
	trainer.EarOn()
	trainer.TeamNames()

	go player(p)

	serverParams := trainer.ServerParams()
	for {
		currentTime := trainer.Time()
		if (currentTime+1)%300 == 0 {
			ballPos := rcsscommon.RandomBallPosition()
			trainer.MoveBall(ballPos.X, ballPos.Y, 0, 0)
		}

		if trainer.PlayMode() == "before_kick_off" {
			trainer.Start()
		}

		err = trainer.Error()
		for err != nil {
			trainer.Log(err)
			err = trainer.Error()
		}

		if serverParams.SynchMode {
			trainer.DoneSynch()
			trainer.WaitSynch()
		} else {
			trainer.WaitNextStep(currentTime)
		}
	}
}
