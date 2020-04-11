package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	stringutils "github.com/alessiosavi/GoGPUtils/string"

	"github.com/alessiosavi/GoStatOgame/datastructure/players"
	"github.com/alessiosavi/GoStatOgame/datastructure/score"
)

const rangePoint = 20 // 20%

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type InputRequest struct {
	Username string `json:"username"`
	Percent  int    `json:"percent"`
	Uni      int    `json:"uni"`
}

func main() {

	var inputRequest InputRequest
	if !stringutils.IsBlank(os.Getenv("console")) {

		username := flag.String("username", "", "Username related to your account")
		percent := flag.Int("percent", -1, "Percent of points for enhance the number of target found")
		uni := flag.Int("uni", -1, "Number of uni related to your account")
		flag.Parse()
		if stringutils.IsBlank(*username) || *percent == -1 || *uni == -1 {
			flag.Usage()
			return
		}

		inputRequest.Username = *username
		inputRequest.Percent = *percent
		inputRequest.Uni = *uni

		// Loading all players from http api
		p, err := players.LoadPlayers(inputRequest.Uni)
		check(err)

		// Find the player related to the input username
		myPlayer, err := p.FilterPlayerByName(inputRequest.Username)
		check(err)

		// Load the data related to the user (planets, alliance, etc)
		myPlayerData, err := players.RetrievePlayerDataByID(inputRequest.Uni, myPlayer.ID)
		check(err)

		// Load the score related to the fleet of the input user
		myScore, err := score.LoadMilitaryScore(inputRequest.Uni)
		check(err)

		// TODO: Populate DynamoDB
		// downloadPlayerData(p)

		toAttack := filterCandidateAttack(myPlayerData, myScore.Player)

		for i := range toAttack {
			player, err := p.FilterPlayerById(toAttack[i].ID)
			check(err)
			toAttack[i].Username = player.Name
		}

		data, _ := json.MarshalIndent(toAttack, "", "  ")
		ioutil.WriteFile("./toAttack.json", data, 0644)
	}

}

// TODO: Retrieve from DynamoDB
func filterCandidateAttack(me players.PlayerData, candidate []score.Player) []players.PlayerData {

	var toAttack []players.PlayerData
	var differenceAllowed int
	for _, c := range candidate {

		candidateData, err := ioutil.ReadFile("/home/alessiosavi/WORKSPACE/Go/GoStatOgame/data/playerdata/" + c.ID + ".xml")
		if err != nil {
			fmt.Printf("ID %s not found\n", c.ID)
			continue
		}
		var playerData players.PlayerData
		if err := xml.Unmarshal(candidateData, &playerData); err != nil {
			fmt.Printf("Panic on ID %s\n", c.ID)
			panic(err)
		}

		playerData.ID = c.ID
		// // Get global position
		// globalPointsTarget, _ := strconv.Atoi(playerData.GetTotal().Score)
		// myGlobalPoints, _ := strconv.Atoi(me.GetTotal().Score)
		// differenceAllowed = percent(rangePoint, myGlobalPoints)

		// if myGlobalPoints+differenceAllowed > globalPointsTarget && globalPointsTarget > myGlobalPoints-differenceAllowed {
		// 	fmt.Printf("You can attack: %+v\n", playerData)
		// 	toAttack = append(toAttack, playerData)
		// }

		// Get fleet position
		fleetPointsTarget, _ := strconv.Atoi(playerData.GetMilitary().Score)
		myFleetGlobalPoints, _ := strconv.Atoi(me.GetMilitary().Score)
		differenceAllowed = percent(rangePoint, myFleetGlobalPoints)

		if myFleetGlobalPoints+differenceAllowed > fleetPointsTarget && fleetPointsTarget > myFleetGlobalPoints-differenceAllowed {
			toAttack = append(toAttack, playerData)
		}

	}
	return toAttack
}

func percent(percent int, all int) int {
	return int((all * percent) / 100)
}
