package actions

import (
	"net/http"
	"net/url"
	"tic-tac-toe__api/models"
	"time"

	"github.com/gobuffalo/buffalo"
)

type EventProcessResult struct {
	ID        int       `json:"id"`
	Payload   string    `json:"payload"`
	EventDt   time.Time `json:"event_dt"`
	ReceiveDt time.Time `json:"receive_dt"`
}

type EventProcessResults []EventProcessResult

type EventProcessResponse struct {
	Results   map[models.EventType]EventProcessResults `json:"results"`
	LastEvent int                                      `json:"last_event"`
}

// EventsGet default implementation.
func EventsGet(c buffalo.Context) error {
	tokens := c.Params().(url.Values)["token[]"]
	lastEvent := c.Param("lastEvent")
	events, err := models.GetLastEvents(tokens, lastEvent)

	if err != nil {
		return c.Render(http.StatusBadRequest, r.JSON(err.Error()))
	}

	if len(events) == 0 {
		// return c.Error(http.StatusNoContent, errors.New("No events"))
		return c.Render(http.StatusNoContent, nil)
	}

	lastEventObject := events[len(events)-1]
	eventsResults := map[models.EventType]EventProcessResults{}

	for _, event := range events {
		payload, err := event.ProcessEventPayload()

		if err != nil {
			continue
			// return c.Render(http.StatusBadRequest, r.JSON(err.Error()))
		}

		eventsResults[event.Type] = append(eventsResults[event.Type], EventProcessResult{
			ID:        event.ID,
			Payload:   payload,
			EventDt:   event.CreatedAt,
			ReceiveDt: time.Now(),
		})
	}

	return c.Render(http.StatusOK, r.JSON(EventProcessResponse{
		Results:   eventsResults,
		LastEvent: lastEventObject.ID,
	}))
}
