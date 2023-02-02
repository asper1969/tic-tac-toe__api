package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

type EventType int

const (
	ROUND_START               EventType = 1
	TEAM_JOIN_TOURNAMENT      EventType = 2
	TEAM_MAKE_MOVE            EventType = 3
	TEAM_ACCEPT_OPPONENT_MOVE EventType = 4
	TEAM_ANSWERED_QUESTION    EventType = 5
	TEAM_PASSED_MOVE          EventType = 6
	TEAM_WINS                 EventType = 7
	TOURNAMENT_PAUSED         EventType = 8
	TOURNAMENT_CONTINUED      EventType = 9
	TOURNAMENT_STOPPED        EventType = 10
	MODERATOR_UPDATES_MATCH   EventType = 11
)

// Event is used by pop to map your events database table to your go code.
type Event struct {
	ID         int       `json:"id" db:"id"`
	SenderID   uuid.UUID `json:"sender_id" db:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id" db:"receiver_id"`
	Type       EventType `json:"type" db:"type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (e Event) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Events is not required by pop and may be deleted
type Events []Event

// String is not required by pop and may be deleted
func (e Events) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (e *Event) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *Event) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *Event) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func GetLastEvents(tokens []string, lastEventID string) (Events, error) {
	events := Events{}
	//get all tokens events
	dbQuery := DB.Where("receiver_id IN (?)", tokens)

	if lastEventID != "" {
		//get all events for tokens after lastEvent create_dt
		lastEvent := Event{}
		err := DB.Where("id = ?", lastEventID).First(&lastEvent)

		if err != nil {
			return nil, err
		}

		dbQuery.Where("created_at > ?", lastEvent.CreatedAt)
	}

	dbQuery.Order("created_at ASC")
	err := dbQuery.All(&events)

	if err != nil {
		return nil, err
	}

	return events, nil
}
