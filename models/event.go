package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

type EventType int

const (
	ROUND_START                 EventType = 1
	TEAM_JOIN_TOURNAMENT        EventType = 2
	TEAM_MAKE_MOVE              EventType = 3
	TEAM_ACCEPT_OPPONENT_MOVE   EventType = 4
	TEAM_ANSWERED_QUESTION      EventType = 5
	TEAM_PASSED_MOVE            EventType = 6
	TEAM_WINS                   EventType = 7
	TOURNAMENT_PAUSED           EventType = 8
	TOURNAMENT_CONTINUED        EventType = 9
	TOURNAMENT_STOPPED          EventType = 10
	MODERATOR_UPDATES_MATCH     EventType = 11
	ROUND_STOPPED               EventType = 12
	TEAM_DECLINED_OPPONENT_MOVE EventType = 13
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

type MeetingResultPayload struct {
	MeetingID  int   `json:"meeting_id"`
	FTeam      *Team `json:"f_team"`
	STeam      *Team `json:"s_team"`
	FTeamScore int   `json:"f_team_score"`
	STeamScore int   `json:"s_team_score"`
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

func (e *Event) ProcessEventPayload() (string, error) {
	var payload []byte

	switch e.Type {
	case ROUND_START:
		//Each team gets their meeting
		tokenID := e.ReceiverID
		meeting, err := GetMeetingByTokenID(tokenID)

		fmt.Println(meeting)

		if err != nil {
			return "", err
		}

		payload, _ = json.Marshal(meeting)

		if err != nil {
			return "", err
		}
	case TEAM_JOIN_TOURNAMENT:
		//Moderator gets team data
		//TODO: all actors gets data
		team, err := GetTeamByTokenID(e.SenderID)

		if err != nil {
			return "", err
		}

		payload, _ = json.Marshal(team)

		if err != nil {
			return "", err
		}
	case TEAM_ACCEPT_OPPONENT_MOVE:
		//Active team get signal answer quiz question
		payload, _ = json.Marshal(map[string]bool{"move_accepted": true})
	case TEAM_MAKE_MOVE, TEAM_ANSWERED_QUESTION:
		//Opponent team gets last meeting_log record
		tokenID := e.ReceiverID
		meetingLog, err := GetLastMeetingLogByTokenID(tokenID)

		if err != nil {
			return "", err
		}

		payload, _ = json.Marshal(meetingLog)

		if err != nil {
			return "", err
		}
	case TEAM_PASSED_MOVE:
		payload, _ = json.Marshal(map[string]bool{"opponent_passed_move": true})
	case TEAM_WINS:
		//All tournament actors gets meeting result
		tokenID := e.SenderID
		meetingLog := MeetingLog{}
		meeting, err := GetMeetingByTokenID(tokenID)

		if err != nil {
			return "", err
		}

		err = DB.Where("meeting_id = ?", meeting.ID).Last(&meetingLog)

		payload, _ = json.Marshal(MeetingResultPayload{
			MeetingID:  meeting.ID,
			FTeam:      meeting.FTeam,
			STeam:      meeting.STeam,
			FTeamScore: meetingLog.FTeamScore,
			STeamScore: meetingLog.STeamScore,
		})

		if err != nil {
			return "", err
		}
	case TOURNAMENT_PAUSED:
		//All teams gets signal. Game freezes
	case TOURNAMENT_CONTINUED:
		//All teams gets signal. Game resumed
	case TOURNAMENT_STOPPED:
		//All teams gets signal. All meetings ends
		//TODO: in payload returns all current meetings results
	case MODERATOR_UPDATES_MATCH:
		//Both meeting teams gets updates
	}

	return string(payload), nil
}

func GetTokenByID(tokenID uuid.UUID) (Token, error) {
	token := Token{}
	err := DB.Where("id = ?", tokenID).Last(&token)
	return token, err
}

func GetTeamByTokenID(tokenID uuid.UUID) (Team, error) {
	team := Team{}

	token, err := GetTokenByID(tokenID)

	if err != nil {
		return team, err
	}

	err = DB.Where("id = ?", token.ObjectID).Last(&team)

	return team, err
}

func GetMeetingByTokenID(tokenID uuid.UUID) (Meeting, error) {
	meeting := Meeting{}
	team, err := GetTeamByTokenID(tokenID)

	if err != nil {
		return meeting, err
	}

	err = DB.Where("f_team_id = ? OR s_team_id = ?", team.ID, team.ID).Last(&meeting)

	if err != nil {
		return meeting, err
	}

	var fTeam, sTeam Team
	err = DB.Where("id = ?", meeting.FTeamID).Last(&fTeam)

	if err != nil {
		return meeting, err
	}

	err = DB.Where("id = ?", meeting.STeamID).First(&sTeam)

	if err != nil {
		return meeting, err
	}

	meeting.FTeam = &fTeam
	meeting.STeam = &sTeam

	// fmt.Println(fTeam)

	return meeting, err
}

func GetLastMeetingLogByTokenID(tokenID uuid.UUID) (MeetingLog, error) {
	meetingLog := MeetingLog{}
	meeting, err := GetMeetingByTokenID(tokenID)

	if err != nil {
		return meetingLog, err
	}

	err = DB.Where("meeting_id = ?", meeting.ID).Last(&meetingLog)

	return meetingLog, err
}
