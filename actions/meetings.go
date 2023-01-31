package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// MeetingsMakeMove default implementation.
func MeetingsMakeMove(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("meetings/make_move.html"))
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

