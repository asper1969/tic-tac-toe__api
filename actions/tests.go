package actions

import (
	"net/http"
	"tic-tac-toe__api/models"

	"github.com/gobuffalo/buffalo"
)

// TestsGet default implementation.
func TestsGet(c buffalo.Context) error {
	res := models.GetMeetingField(1, 6, 3, 5)
	return c.Render(http.StatusOK, r.JSON(res))
}

// TestsPost default implementation.
func TestsPost(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("tests/post.html"))
}
