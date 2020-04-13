package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alessiosavi/GoStatOgame/datastructure/players"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"os"
	"strconv"

	stringutils "github.com/alessiosavi/GoGPUtils/string"

	"github.com/alessiosavi/GoStatOgame/datastructure/score"
)

const rangePoint = 30 // 20%

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

		toAttack, _ := core(inputRequest)

		data, _ := json.MarshalIndent(toAttack, "", "  ")
		ioutil.WriteFile("toAttack.json", data, 0644)
		fmt.Println(string(data))
	} else {
		lambda.Start(core)
	}

}

func core(inputRequest InputRequest) ([]score.Player, error) {
	// Loading all players from http api
	p, err := players.LoadPlayers(inputRequest.Uni)
	check(err)

	// Find the player related to the input username
	myPlayer, err := p.FilterPlayerByName(inputRequest.Username)
	check(err)

	// Load the score related to the fleet of all the users
	allMilitaryScore, err := score.LoadMilitaryScore(inputRequest.Uni)
	check(err)

	myScore, err := allMilitaryScore.FilterScoreByID(myPlayer.ID)
	check(err)

	toAttack := filterCandidateAttack(myScore, allMilitaryScore.Player)

	for i := range toAttack {
		player, err := p.FilterPlayerById(toAttack[i].ID)
		check(err)
		toAttack[i].Username = player.Name
	}
	return toAttack, nil
}

// TODO: Retrieve from DynamoDB
func filterCandidateAttack(me score.Player, candidate []score.Player) []score.Player {

	var toAttack []score.Player
	var differenceAllowed int
	for _, c := range candidate {

		// Get fleet position
		fleetPointsTarget, _ := strconv.Atoi(c.Score)
		myFleetGlobalPoints, _ := strconv.Atoi(me.Score)
		differenceAllowed = percent(rangePoint, myFleetGlobalPoints)

		if myFleetGlobalPoints+differenceAllowed > fleetPointsTarget && fleetPointsTarget > myFleetGlobalPoints-differenceAllowed {
			toAttack = append(toAttack, c)
		}

	}
	return toAttack
}

func percent(percent int, all int) int {
	return (all * percent) / 100
}
