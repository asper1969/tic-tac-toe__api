package actions

import (
	"fmt"
	"net/http"
	"tic-tac-toe__api/models"

	"github.com/gobuffalo/buffalo"
)

type MoveRequestPayload struct {
	TokenTeam    string `json:"token_team"`
	PlacesSet    string `json:"places_set"`
	QuestionsLog string `json:"questions_log"`
	FTeamScore   int    `json:"f_team_score"`
	STeamScore   int    `json:"s_team_score"`
}

type TeamWinRequestPayload struct {
	TokenTeam string `json:"token_team"`
}

// MeetingsMakeMove default implementation.
//TODO: test
func MeetingsMakeMove(c buffalo.Context) error {
	payload := &models.TeamActionPayload{}

	if err := c.Bind(payload); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	err := models.ProcessTeamAction(models.TEAM_MAKE_MOVE, *payload)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Move was saved"}))
}

// MeetingsAnswerQuestion default implementation.
func MeetingsAnswerQuestion(c buffalo.Context) error {
	payload := &models.TeamActionPayload{}

	if err := c.Bind(payload); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	err := models.ProcessTeamAction(models.TEAM_ANSWERED_QUESTION, *payload)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Answer was saved"}))
}

// MeetingsAcceptMove default implementation.
func MeetingsAcceptMove(c buffalo.Context) error {
	payload := &models.TeamActionPayload{}

	if err := c.Bind(payload); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	err := models.ProcessTeamAction(models.TEAM_ACCEPT_OPPONENT_MOVE, *payload)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Move was accepted"}))
}

// MeetingsPassMove default implementation.
func MeetingsPassMove(c buffalo.Context) error {
	payload := &models.TeamActionPayload{}

	if err := c.Bind(payload); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	err := models.ProcessTeamAction(models.TEAM_PASSED_MOVE, *payload)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Move was passed"}))
}

// MeetingsTeamWin default implementation.
func MeetingsTeamWin(c buffalo.Context) error {
	payload := &TeamWinRequestPayload{}

	if err := c.Bind(payload); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	err := models.ProcessTeamWin(payload.TokenTeam)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Your team win!"}))
}
