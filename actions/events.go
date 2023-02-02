package actions

import (
	"net/http"
	"net/url"
	"tic-tac-toe__api/models"
	"time"

	"github.com/gobuffalo/buffalo"
)

type EventProcessResult struct {
	ID        int               `json:"id"`
	Payload   map[string]string `json:"payload"`
	EventDt   time.Time         `json:"event_dt"`
	ReceiveDt time.Time         `json:"receive_dt"`
}

type EventProcessResults []EventProcessResult

// EventsGet default implementation.
func EventsGet(c buffalo.Context) error {
	tokens := c.Params().(url.Values)["token[]"]
	lastEvent := c.Param("lastEvent")
	events, err := models.GetLastEvents(tokens, lastEvent)

	if err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(err.Error()))
	}

	eventsResults := map[models.EventType]EventProcessResults{}

	for _, event := range events {
		payload, err := event.ProcessEventPayload()

		if err != nil {
			return c.Render(http.StatusBadRequest, r.JSON(err.Error()))
		}

		eventsResults[event.Type] = append(eventsResults[event.Type], EventProcessResult{
			ID:        event.ID,
			Payload:   payload,
			EventDt:   event.CreatedAt,
			ReceiveDt: time.Now(),
		})
	}

	return c.Render(http.StatusOK, r.JSON(eventsResults))
}
