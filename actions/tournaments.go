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
	Locale       string `json:"lang"`
	MaxScore     int    `json:"maxScore"`
	TeamsAmount  int    `json:"numberOfTeams"`
	FieldsAmount int    `json:"numberOfFields"`
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

type JoinTournamentResponse struct {
	TokenTeam        string `json:"token_team"`
	TokenTournament  string `json:"token_tournament"`
	SettingsMaxScore int    `json:"settings:max_score"`
	SettingsLang     string `json:"settings:lang"`
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
		GamePass:     gameCode,
		Locale:       requestData.Locale,
		MaxScore:     requestData.MaxScore,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		TeamsAmount:  requestData.TeamsAmount,
		FieldsAmount: requestData.FieldsAmount,
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
		Type:      models.TOKEN_TOURNAMENT, //tournament
		ObjectID:  tournament.ID,
		CreatedAt: time.Now(),
	}
	tokenModerator := models.Token{
		ID:        moderatorUUID,
		Type:      models.TOKEN_MODERATOR, //moderator
		ObjectID:  tournament.ID,
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

	//Get tournament by code
	tournament := models.Tournament{}

	err := models.DB.Where("game_pass = ?", requestData.Code).First(&tournament)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Create team record
	team := models.Team{
		Name:         requestData.TeamName,
		TournamentID: tournament.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = models.DB.Create(&team)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	//Create team token
	teamUUID, err := uuid.NewV7()

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	tokenTeam := models.Token{
		ID:        teamUUID,
		Type:      models.TOKEN_TEAM, //team
		ObjectID:  team.ID,
		CreatedAt: time.Now(),
	}

	err = models.DB.Create(&tokenTeam)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	tournamentToken := tournament.GetToken().ID

	//Get moderator token
	//Create new event
	event := models.Event{
		SenderID:   tokenTeam.ID,
		ReceiverID: tournamentToken,
		Type:       models.TEAM_JOIN_TOURNAMENT,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = models.DB.Create(&event)

	if err != nil {
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	fmt.Println(tournament)
	//return team token and tournament token

	return c.Render(http.StatusOK, r.JSON(JoinTournamentResponse{
		TokenTeam:        tokenTeam.ID.String(),
		TokenTournament:  tournamentToken.String(),
		SettingsMaxScore: tournament.MaxScore,
		SettingsLang:     tournament.Locale,
	}))
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

// Create new tournament round default implementation.
func TournamentsCreateRound(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Get tournament by token
	tournament := models.Tournament{}
	q := models.DB.Q()
	q.LeftJoin("tokens", "tokens.object_id = tournaments.id")
	q.Where(`tokens.id = ?`, requestData.Token)
	err := q.Last(&tournament)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//update tournament (update start dt)(?)
	//create matches for all team pairs
	err = tournament.CreateNextRound()

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get all tournament matches with start_dt != NULL and end_dt == NULL
	//In response return all matches
	meetings := models.Meetings{}
	err = models.DB.Where("start_dt IS NULL AND end_dt IS NULL AND tournament_id = ?", tournament.ID).All(&meetings)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	for i, meeting := range meetings {
		teams := models.Teams{}
		err := models.DB.Where("id = ? || id = ?", meeting.FTeamID, meeting.STeamID).All(&teams)

		if err != nil {
			return c.Render(http.StatusOK, r.JSON(err.Error()))
		}

		meetings[i].FTeam = &teams[0]
		meetings[i].STeam = &teams[1]
	}

	return c.Render(http.StatusOK, r.JSON(meetings))
}

// Start tournament round default implementation.
func TournamentsStartRound(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Get tournament by token
	tournament := models.Tournament{}
	q := models.DB.Q()
	q.LeftJoin("tokens", "tokens.object_id = tournaments.id")
	q.Where(`tokens.id = ?`, requestData.Token)
	err := q.Last(&tournament)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Start first round
	meetings, err := tournament.StartNextRound()

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(meetings))
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
