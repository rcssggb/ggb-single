package player

import (
	"sync"

	"github.com/rcssggb/ggb-lib/playerclient"
)

// Player is the high-level structure containing player methods and sensors
type Player struct {
	Client *playerclient.Client

	mutex              sync.RWMutex
	friendlyPlayersPos map[int]SeenPlayerPosition
	opponentPlayersPos map[int]SeenPlayerPosition
	self               SelfData
	ball               Ball
}

// NewPlayer constructs and initializes the player struct
func NewPlayer(team, host string) (*Player, error) {
	client, err := playerclient.NewPlayerClient(team, host)
	if err != nil {
		client.Bye()
		return nil, err
	}

	p := &Player{
		Client: client,
	}

	go p.bodyUpdate()
	go p.sightUpdate()

	return p, nil
}
