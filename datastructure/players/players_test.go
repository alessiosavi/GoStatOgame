package players

import (
	"testing"
)

const dataPath = "/home/alessiosavi/WORKSPACE/Go/GoStatOgame/data/players.xml"

func Test_LoadPlayer(t *testing.T) {
	players, err := LoadPlayers(dataPath)
	if err != nil {
		t.Error(err)
	}
	t.Log("Players len: ", len(players.Players))
}

func Test_FilterPlayerByName(t *testing.T) {
	var player Player
	var err error
	players, err := LoadPlayers(dataPath)
	if err != nil {
		t.Error(err)
	}

	if player, err = players.FilterPlayerByName("Geologist Lepus"); err != nil {
		t.Error(err)
	}
	t.Logf("Player: %+v\n", player)
}
