package main

import (
	"log"
	"time"

	"github.com/rcssggb/ggb-lib/playerclient"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	hostName := "rcssserver"

	for {
		p, err := playerclient.NewPlayerClient("single-agent", hostName)
		if err != nil {
			log.Println(err)
			continue
		}

		time.Sleep(2 * time.Second)

		player(p)
	}
}
