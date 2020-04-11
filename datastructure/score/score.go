package score

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/alessiosavi/GoStatOgame/utils"
)

type Player struct {
	Position string `xml:"position,attr"`
	ID       string `xml:"id,attr"`
	Score    string `xml:"score,attr"`
	Ships    string `xml:"ships,attr"`
}

type MilitaryScore struct {
	Player []Player `xml:"player"`
}

type MilitaryScoreLost MilitaryScore

func LoadMilitaryScore(uni int) (MilitaryScore, error) {
	var players MilitaryScore
	var err error
	var resp *http.Response
	var body []byte
	var url = fmt.Sprintf("https://s%d-it.ogame.gameforge.com/api/highscore.xml?category=1&type=3", uni)

	if resp, err = http.Get(url); err != nil {
		return players, err
	}
	if body, err = utils.ReadBody(resp); err != nil {
		return players, err
	}
	err = xml.Unmarshal(body, &players)
	return players, err
}

func LoadMilitaryScoreLost(url string) (MilitaryScoreLost, error) {
	var players MilitaryScoreLost
	var data []byte
	var err error
	// TODO: Change with http request
	if data, err = ioutil.ReadFile(url); err != nil {
		return MilitaryScoreLost{}, err
	}

	err = xml.Unmarshal(data, &players)
	return players, err
}

func (s *MilitaryScore) FilterLowerPlayerById(id string) ([]Player, error) {
	var players []Player

	for i := range s.Player {
		if strings.Compare(s.Player[i].ID, id) == 0 {
			return s.Player[i:], nil
		}
	}

	if len(players) > 0 {
		return players, nil
	}
	return nil, fmt.Errorf("players below [%s] not found", id)
}

func (s *MilitaryScore) FilterHigherPlayerById(id string) ([]Player, error) {
	for i := range s.Player {
		if strings.Compare(s.Player[i].ID, id) == 0 {
			return s.Player[:i], nil
		}
	}
	return nil, fmt.Errorf("players higher than [%s] not found", id)
}
