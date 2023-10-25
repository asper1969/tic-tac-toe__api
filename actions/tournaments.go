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

type JoinOpponentsRequest struct {
	FTeamName string `json:"f_team_name"`
	STeamName string `json:"s_team_name"`
	Code      string `json:"code"`
}

type JoinOpponentsResponse struct {
	FTeamName        string `json:"f_team_name"`
	FTeamToken       string `json:"f_team_token"`
	STeamName        string `json:"s_team_name"`
	STeamToken       string `json:"s_team_token"`
	TokenTournament  string `json:"token_tournament"`
	SettingsMaxScore int    `json:"settings:max_score"`
	SettingsLang     string `json:"settings:lang"`
	SettingsField    int    `json:"settings:field"`
	Field            int    `json:"field"`
	ModeratorToken   string `json:"moderator_token"`
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
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	//Generate tournament token and moderator token
	//Save tokens in db
	tournamenUUID, err := uuid.NewV7()

	if err != nil {
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	moderatorUUID, err := uuid.NewV7()

	if err != nil {
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	tokenTournament := models.Token{
		ID:        tournamenUUID,
		Type:      models.TOKEN_TOURNAMENT, //tournament
		ObjectID:  tournament.ID,
		CreatedAt: time.Now(),
		Expired:   false,
	}
	tokenModerator := models.Token{
		ID:        moderatorUUID,
		Type:      models.TOKEN_MODERATOR, //moderator
		ObjectID:  tournament.ID,
		CreatedAt: time.Now(),
		Expired:   false,
	}

	err = models.DB.Create(&tokenTournament)

	if err != nil {
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	err = models.DB.Create(&tokenModerator)

	if err != nil {
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
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Get tournament by code
	tournament := models.Tournament{}

	err := models.DB.Where("game_pass = ?", requestData.Code).First(&tournament)

	if err != nil {
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
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	//Create team token
	teamUUID, err := uuid.NewV7()

	if err != nil {
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	tokenTeam := models.Token{
		ID:        teamUUID,
		Type:      models.TOKEN_TEAM, //team
		ObjectID:  team.ID,
		CreatedAt: time.Now(),
		Expired:   false,
	}

	err = models.DB.Create(&tokenTeam)

	if err != nil {
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
		return c.Render(http.StatusBadGateway, r.JSON(err.Error()))
	}

	//return team token and tournament token

	return c.Render(http.StatusOK, r.JSON(JoinTournamentResponse{
		TokenTeam:        tokenTeam.ID.String(),
		TokenTournament:  tournamentToken.String(),
		SettingsMaxScore: tournament.MaxScore,
		SettingsLang:     tournament.Locale,
	}))
}

// TournamentsJoin default implementation.
func TournamentsJoinOpponents(c buffalo.Context) error {
	requestData := &JoinOpponentsRequest{}

	if err := c.Bind(requestData); err != nil {
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Get tournament by code
	tournament := models.Tournament{}

	err := models.DB.Where("game_pass = ?", requestData.Code).Last(&tournament)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err))
	}

	fTeam, fTeamToken, err := tournament.AddNewTeam(requestData.FTeamName)

	if err != nil {
		return c.Render(http.StatusBadGateway, r.JSON(err.Error()))
	}

	sTeam, sTeamToken, err := tournament.AddNewTeam(requestData.FTeamName)

	if err != nil {
		return c.Render(http.StatusBadGateway, r.JSON(err.Error()))
	}

	moderatorToken := tournament.GetToken()

	//return team token and tournament token
	return c.Render(http.StatusOK, r.JSON(JoinOpponentsResponse{
		FTeamName:        fTeam.Name,
		FTeamToken:       fTeamToken.ID.String(),
		STeamName:        sTeam.Name,
		STeamToken:       sTeamToken.ID.String(),
		TokenTournament:  tournament.GetToken().String(),
		SettingsMaxScore: tournament.MaxScore,
		SettingsLang:     tournament.Locale,
		SettingsField:    tournament.FieldsAmount,
		Field:            moderatorToken.ObjectID,
		ModeratorToken:   moderatorToken.ID.String(),
	}))
}

// TournamentsAction default implementation.
func TournamentsAction(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(requestData))
}

// Create new tournament round default implementation.
func TournamentsCreateRound(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get tournament by token
	tournament := models.Tournament{}
	q := models.DB.Q()
	q.LeftJoin("tokens", "tokens.object_id = tournaments.id")
	q.Where(`tokens.id = ?`, requestData.Token)
	err := q.Last(&tournament)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
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
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	for i, meeting := range meetings {
		teams := models.Teams{}
		err := models.DB.Where("id = ? OR id = ?", meeting.FTeamID, meeting.STeamID).All(&teams)

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
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Get tournament by token
	tournament := models.Tournament{}
	q := models.DB.Q()
	q.LeftJoin("tokens", "tokens.object_id = tournaments.id")
	q.Where(`tokens.id = ?`, requestData.Token)
	err := q.Last(&tournament)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Start first round
	meetings, err := tournament.StartNextRound()

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(meetings))
}

// Start tournament round default implementation.
func TournamentsStopRound(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	token := models.Token{}
	err := models.DB.Where("id = ?", requestData.Token).Last(&token)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get tournament by tournamentToken
	tournament := models.Tournament{}
	err = models.DB.Where("id = ?", token.ObjectID).Last(&tournament)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get all tournament meetings
	meetings := models.Meetings{}
	err = models.DB.Where("tournament_id = ? AND end_dt IS NULL", tournament.ID).All(&meetings)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	for _, meeting := range meetings {
		models.EndMeeting(meeting.ID)
	}

	event := models.Event{
		SenderID:   token.ID,
		ReceiverID: token.ID,
		Type:       models.ROUND_STOPPED,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = models.DB.Create(&event)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON("round has been finished"))
}

// TournamentsStop default implementation.
func TournamentsStop(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	token := models.Token{}
	err := models.DB.Where("id = ?", requestData.Token).Last(&token)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get tournament by tournamentToken
	tournament := models.Tournament{}
	err = models.DB.Where("id = ?", token.ObjectID).Last(&tournament)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get all tournament meetings
	meetings := models.Meetings{}
	err = models.DB.Where("tournament_id = ? AND end_dt IS NULL", tournament.ID).All(&meetings)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	for _, meeting := range meetings {
		models.EndMeeting(meeting.ID)
	}

	event := models.Event{
		SenderID:   token.ID,
		ReceiverID: token.ID,
		Type:       models.TOURNAMENT_STOPPED,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = models.DB.Create(&event)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON("tournament has been finished"))
}

// TournamentsPause default implementation.
func TournamentsPause(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	token := models.Token{}
	err := models.DB.Where("id = ?", requestData.Token).Last(&token)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//TODO: create event from tournament_token to
	event := models.Event{
		SenderID:   token.ID,
		ReceiverID: token.ID,
		Type:       models.TOURNAMENT_PAUSED,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = models.DB.Create(&event)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON("tournament has been paused"))
}

// TournamentsPause default implementation.
func TournamentsContinue(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	token := models.Token{}
	err := models.DB.Where("id = ?", requestData.Token).Last(&token)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//TODO: create event from tournament_token to
	event := models.Event{
		SenderID:   token.ID,
		ReceiverID: token.ID,
		Type:       models.TOURNAMENT_CONTINUED,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = models.DB.Create(&event)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON("tournament has been continued"))
}

// TournamentsPause default implementation.
func TournamentsEnd(c buffalo.Context) error {
	requestData := &ActionTournamentRequest{}

	if err := c.Bind(requestData); err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	token := models.Token{}
	err := models.DB.Where("id = ?", requestData.Token).Last(&token)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//TODO: create event from tournament_token to
	event := models.Event{
		SenderID:   token.ID,
		ReceiverID: token.ID,
		Type:       models.TOURNAMENT_ENDED,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = models.DB.Create(&event)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON("tournament has been paused"))
}

func TournamentTokenIsActive(c buffalo.Context) error {
	tournamentTokenId := c.Param("tournamentToken")

	if tournamentTokenId == "" {
		return c.Render(http.StatusBadRequest, r.JSON("tournamentToken is not specified"))
	}

	token := models.Token{}

	err := models.DB.Where("id = ? AND expired = FALSE", tournamentTokenId).Last(&token)

	fmt.Println(token)

	if err != nil {
		return c.Render(http.StatusNoContent, nil)
	}

	return c.Render(http.StatusOK, r.JSON(token))
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

func TournamentGetSession(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(""))
}
