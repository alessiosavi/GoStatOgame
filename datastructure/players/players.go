package players

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/alessiosavi/OgameStats/utils"
)

type Player struct {
	ID       string `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	Status   string `xml:"status,attr"`
	Alliance string `xml:"alliance,attr"`
}

type Players struct {
	Players []Player `xml:"player"`
}

// PlayerData.Position.Type
// 0 - Total
// 1 - Economy
// 2 - Research
// 3 - Military
// 4 - Military Lost
// 5 - Military Built
// 6 - Military Destoyed
// 7 - Honor
type PlayerData struct {
	Username  string
	ID        string
	Positions struct {
		Position []Position `xml:"position"`
	} `xml:"positions"`
	Planets struct {
		Planet []Planet `xml:"planet"`
	} `xml:"planets"`
	Alliance struct {
		ID   string `xml:"id,attr"`
		Name string `xml:"name"`
		Tag  string `xml:"tag"`
	} `xml:"alliance"`
}

type Position struct {
	Type  string `xml:"type,attr"`
	Score string `xml:"score,attr"`
	Ships string `xml:"ships,attr"`
}

type Planet struct {
	ID     string `xml:"id,attr"`
	Name   string `xml:"name,attr"`
	Coords string `xml:"coords,attr"`
}

func LoadPlayers(uni int) (Players, error) {
	var players Players
	var err error
	var resp *http.Response
	var body []byte
	var url = fmt.Sprintf("https://s%d-it.ogame.gameforge.com/api/players.xml", uni)

	if resp, err = http.Get(url); err != nil {
		return players, err
	}
	body, _ = utils.ReadBody(resp)
	err = xml.Unmarshal(body, &players)
	return players, err
}

// func (p *Player) GetPlayerData(url string) (PlayerData, error) {
// 	var playerData PlayerData
// 	var data []byte
// 	var err error
// 	// TODO: Change with http request
// 	if data, err = ioutil.ReadFile(url); err != nil {
// 		return PlayerData{}, err
// 	}
// 	if err = xml.Unmarshal(data, &playerData); err != nil {
// 		return PlayerData{}, err
// 	}
// 	return playerData, nil
// }

func RetrievePlayerDataByID(uni int, id string) (PlayerData, error) {
	// data, _ := ioutil.ReadFile("/home/alessiosavi/WORKSPACE/Go/GoStatOgame/data/playerdata/" + id + ".xml")
	var playerData PlayerData
	var err error
	var resp *http.Response
	var body []byte

	var url = fmt.Sprintf("https://s%d-it.ogame.gameforge.com/api/playerData.xml?id=%s", uni, id)
	if resp, err = http.Get(url); err != nil {
		return playerData, err
	}
	if body, err = utils.ReadBody(resp); err != nil {
		return playerData, err
	}

	if err := xml.Unmarshal(body, &playerData); err != nil {
		return playerData, err
	}
	return playerData, nil
}

func (p *Players) FilterPlayerByName(username string) (Player, error) {
	for i := range p.Players {
		if strings.EqualFold(p.Players[i].Name, username) {
			return p.Players[i], nil
		}
	}
	return Player{}, fmt.Errorf("player with name [%s] not found", username)
}

func (p *Players) FilterPlayerById(id string) (Player, error) {
	for i := range p.Players {
		if strings.Compare(p.Players[i].ID, id) == 0 {
			return p.Players[i], nil
		}
	}
	return Player{}, fmt.Errorf("player with ID [%s] not found", id)
}

// Wrapper method for extract the position data
func (p *PlayerData) GetTotal() Position {
	return p.Positions.Position[0]
}

func (p *PlayerData) GetEconomy() Position {
	return p.Positions.Position[1]
}

func (p *PlayerData) GetResearch() Position {
	return p.Positions.Position[2]
}

func (p *PlayerData) GetMilitary() Position {
	return p.Positions.Position[3]
}

func (p *PlayerData) GetMilitaryLost() Position {
	return p.Positions.Position[4]
}

func (p *PlayerData) GetMilitaryBuilt() Position {
	return p.Positions.Position[5]
}

func (p *PlayerData) GetMilitaryDestoyed() Position {
	return p.Positions.Position[6]
}
func (p *PlayerData) GetHonor() Position {
	return p.Positions.Position[7]
}
