package actions

import (
	"fmt"
	"net/http"
	"tic-tac-toe__api/models"

	"github.com/gobuffalo/buffalo"
	"github.com/thanhpk/randstr"
)

type CreateTournamentRequest struct {
	MaxScore int `json:"max_score"`
	Rounds   int `json:"rounds"`
}

type JoinTournamentRequest struct {
	TeamName string `json:"team_name"`
	Code     string `json:"code"`
}

type ActionTournamentRequest struct {
	Token  string `json:"token"`
	Action string `json:"action"`
}

// TournamentsCreate default implementation.
func TournamentsCreate(c buffalo.Context) error {
	requestData := &CreateTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//TODO:
	//Create

	return c.Render(http.StatusOK, r.JSON(requestData))
}

// TournamentsJoin default implementation.
func TournamentsJoin(c buffalo.Context) error {
	requestData := &JoinTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(requestData))
}

// TournamentsAction default implementation.
func TournamentsAction(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(requestData))
}

func generateTournamentPass() string {
	str := randstr.String(5)
	tournament := models.Tournament{}
	err := models.DB.Where("game_pass = ?", str).First(&tournament)

	if err == nil {
		str = generateGamePass()
	}

	return str
}
