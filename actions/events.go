package actions

import (
	"net/http"
	"net/url"
	"tic-tac-toe__api/models"

	"github.com/gobuffalo/buffalo"
)

// EventsGet default implementation.
func EventsGet(c buffalo.Context) error {
	tokens := c.Params().(url.Values)["token[]"]
	lastEvent := c.Param("lastEvent")
	events, err := models.GetLastEvents(tokens, lastEvent)

	if err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(events))
}
