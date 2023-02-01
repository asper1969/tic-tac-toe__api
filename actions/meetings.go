package actions

import (
	"net/http"
	"tic-tac-toe__api/models"
	"time"

	"github.com/gobuffalo/buffalo"
)

type MoveRequestPayload struct {
	TokenTeam    string `json:"token_team"`
	PlacesSet    string `json:"places_set"`
	QuestionsLog string `json:"questions_log"`
	FTeamScore   int    `json:"f_team_score"`
	STeamScore   int    `json:"s_team_score"`
}

// MeetingsMakeMove default implementation.
//TODO: test
func MeetingsMakeMove(c buffalo.Context) error {
	payload := &MoveRequestPayload{}

	if err := c.Bind(payload); err != nil {
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Get meeting by team token

	//Get token
	token := &models.Token{}
	err := models.DB.Where("id = ?", payload.TokenTeam).Last(&token)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get team
	team := &models.Team{}
	err = models.DB.Where("id = ?", token.ObjectID).Last(&team)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get meeting
	meeting := &models.Meeting{}
	err = models.DB.Where("f_team_id = ? OR s_team_id = ?", team.ID, team.ID).Last(&meeting)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Create new meeting_log record
	meetingLogRecord := &models.MeetingLog{
		MeetingID:    meeting.ID,
		PlacesSet:    payload.PlacesSet,
		QuestionsLog: payload.QuestionsLog,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ActiveTeam:   team.ID,
		FTeamScore:   payload.FTeamScore,
		STeamScore:   payload.STeamScore,
	}

	err = models.DB.Create(&meetingLogRecord)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get opponent team
	opponent := &models.Team{}
	q := models.DB.Q()

	if team.ID == meeting.FTeamID {
		q.Where("id = ?", meeting.STeamID)
	} else {
		q.Where("id = ?", meeting.FTeamID)
	}

	err = q.Last(&opponent)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Get opponent token
	opponentToken := &models.Token{}
	err = models.DB.Where("object_id = ?", opponent.ID).Last(&opponentToken)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	//Create new TEAM_MAKE_MOVE event
	event := &models.Event{
		SenderID:   token.ID,
		ReceiverID: opponentToken.ID,
		Type:       models.TEAM_MAKE_MOVE,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = models.DB.Create(&event)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(event))
}

// MeetingsAnswerQuestion default implementation.
func MeetingsAnswerQuestion(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("meetings/answer_question.html"))
}

// MeetingsAcceptMove default implementation.
func MeetingsAcceptMove(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("meetings/accept_move.html"))
}

// MeetingsPassMove default implementation.
func MeetingsPassMove(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("meetings/pass_move.html"))
}
