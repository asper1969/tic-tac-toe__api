package actions

import (
	"fmt"
	"net/http"
	"tic-tac-toe__api/models"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
	"github.com/thanhpk/randstr"
)

type CreateTournamentRequest struct {
	Locale   string `json:"locale"`
	MaxScore int    `json:"max_score"`
	Rounds   int    `json:"rounds"`
}

type CreateTournamentResponse struct {
	Code            string `json:"code"`
	TokenModerator  string `json:"token_moderator"`
	TokenTournament string `json:"token_tournament"`
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
	//Generate code

	gameCode := generateTournamentCode()

	//Create new tournament record
	tournament := models.Tournament{
		GamePass:  gameCode,
		Locale:    requestData.Locale,
		MaxScore:  requestData.MaxScore,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Rounds:    requestData.Rounds,
	}

	err := models.DB.Create(&tournament)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	//Generate tournament token and moderator token
	//Save tokens in db
	tournamenUUID, err := uuid.NewV7()

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	moderatorUUID, err := uuid.NewV7()

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	tokenTournament := models.Token{
		ID:        tournamenUUID,
		Type:      1, //tournament
		CreatedAt: time.Now(),
	}
	tokenModerator := models.Token{
		ID:        moderatorUUID,
		Type:      2, //moderator
		CreatedAt: time.Now(),
	}

	err = models.DB.Create(&tokenTournament)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	err = models.DB.Create(&tokenModerator)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	//in response return tokens and tournament code

	return c.Render(http.StatusOK, r.JSON(CreateTournamentResponse{
		Code:            gameCode,
		TokenModerator:  tokenModerator.ID.String(),
		TokenTournament: tokenTournament.ID.String(),
	}))
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

// TournamentsStart default implementation.
func TournamentsStart(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(requestData))
}

// TournamentsStop default implementation.
func TournamentsStop(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(requestData))
}

// TournamentsPause default implementation.
func TournamentsPause(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(requestData))
}

func generateTournamentCode() string {
	str := randstr.String(5)
	tournament := models.Tournament{}
	err := models.DB.Where("game_pass = ?", str).First(&tournament)

	if err == nil {
		str = generateTournamentCode()
	}

	return str
}
