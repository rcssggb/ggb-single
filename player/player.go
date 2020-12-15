package player

import (
	"github.com/rcssggb/ggb-lib/playerclient"
)

// Player is the high-level structure containing player methods and sensors
type Player struct {
	Client  *playerclient.Client
	SelfPos Position
}

// NewPlayer constructs and initializes the player struct
func NewPlayer(team, host string) (*Player, error) {
	client, err := playerclient.NewPlayerClient(team, host)
	if err != nil {
		client.Bye()
		return nil, err
	}

	return &Player{
		Client: client,
	}, nil
}
