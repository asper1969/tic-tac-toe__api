package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// Meeting is used by pop to map your meetings database table to your go code.
type Meeting struct {
	ID           int          `json:"id" db:"id"`
	FTeamID      int          `json:"f_team_id" db:"f_team_id"`
	STeamID      int          `json:"s_team_id" db:"s_team_id"`
	StartDt      nulls.Time   `json:"start_dt" db:"start_dt"`
	EndDt        nulls.Time   `json:"end_dt" db:"end_dt"`
	TournamentID int          `json:"tournament_id" db:"tournament_id"`
	QuestionsSet string       `json:"questions_set" db:"questions_set"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	MeetingLogs  []MeetingLog `json:"meeting_logs" has_many:"meeting_logs"`
	Round        int          `json:"round" db:"round"`
	Field        int          `json:"field" db:"field"`
	FTeam        *Team        `json:"f_team" has_one:"team" fk_id:"f_team_id"`
	STeam        *Team        `json:"s_team" has_one:"team" fk_id:"s_team_id"`
}

type TeamActionPayload struct {
	TokenTeam    string `json:"token_team"`
	PlacesSet    string `json:"places_set"`
	QuestionsLog string `json:"questions_log"`
	FTeamScore   int    `json:"f_team_score"`
	STeamScore   int    `json:"s_team_score"`
}

// String is not required by pop and may be deleted
func (m Meeting) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Meetings is not required by pop and may be deleted
type Meetings []Meeting

// String is not required by pop and may be deleted
func (m Meetings) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Meeting) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Meeting) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Meeting) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func ProcessTeamAction(action EventType, payload TeamActionPayload) error {

	//Get token
	token := Token{}
	err := DB.Where("id = ?", payload.TokenTeam).Last(&token)

	if err != nil {
		return err
	}

	//Get team
	team := Team{}
	err = DB.Where("id = ?", token.ObjectID).Last(&team)

	if err != nil {
		return err
	}

	//Get meeting
	meeting := Meeting{}
	err = DB.Where("f_team_id = ? OR s_team_id = ?", team.ID, team.ID).Last(&meeting)

	if err != nil {
		return err
	}

	//Get opponent team
	opponent := Team{}
	q := DB.Q()

	if team.ID == meeting.FTeamID {
		q.Where("id = ?", meeting.STeamID)
	} else {
		q.Where("id = ?", meeting.FTeamID)
	}

	err = q.Last(&opponent)

	if err != nil {
		return err
	}

	//Get opponent token
	opponentToken := Token{}
	err = DB.Where("object_id = ?", opponent.ID).Last(&opponentToken)

	if err != nil {
		return err
	}

	if action == TEAM_MAKE_MOVE || action == TEAM_ANSWERED_QUESTION || action == TEAM_PASSED_MOVE {
		//Create new meeting_log record
		meetingLogRecord := MeetingLog{
			MeetingID:    meeting.ID,
			PlacesSet:    payload.PlacesSet,
			QuestionsLog: payload.QuestionsLog,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			ActiveTeam:   team.ID,
			FTeamScore:   payload.FTeamScore,
			STeamScore:   payload.STeamScore,
			Accepted:     action == TEAM_ANSWERED_QUESTION,
		}

		err = DB.Create(&meetingLogRecord)

		if err != nil {
			return err
		}
	}

	if action == TEAM_DECLINED_OPPONENT_MOVE {
		//TODO: get meetingLog record before the last one
		lastMeetingLogRecords := []MeetingLog{}
		err = DB.Where("meeting_id = ?", meeting.ID).Order("created_at DESC").Limit(2).All(&lastMeetingLogRecords)

		if err != nil {
			return err
		}

		meetingLogRecord := MeetingLog{
			MeetingID: meeting.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Accepted:  true,
		}

		//Only one record
		if len(lastMeetingLogRecords) < 2 {
			meetingLogRecord.PlacesSet = "[0,0,0,0,0,0,0,0,0]"
			meetingLogRecord.QuestionsLog = "[]"
			meetingLogRecord.ActiveTeam = opponent.ID
			meetingLogRecord.FTeamScore = 0
			meetingLogRecord.STeamScore = 0
		} else {
			previousMeetingLogRecord := lastMeetingLogRecords[1]
			meetingLogRecord.PlacesSet = previousMeetingLogRecord.PlacesSet
			meetingLogRecord.QuestionsLog = previousMeetingLogRecord.QuestionsLog
			meetingLogRecord.ActiveTeam = previousMeetingLogRecord.ActiveTeam
			meetingLogRecord.FTeamScore = previousMeetingLogRecord.FTeamScore
			meetingLogRecord.STeamScore = previousMeetingLogRecord.STeamScore
		}

		//Create new meetingLog
		DB.Create(&meetingLogRecord)

		//Create an event for sender
		event := Event{
			SenderID:   token.ID,
			ReceiverID: token.ID,
			Type:       action,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err = DB.Create(&event)

		if err != nil {
			return err
		}
	}

	//Create new event for opponent
	event := Event{
		SenderID:   token.ID,
		ReceiverID: opponentToken.ID,
		Type:       action,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = DB.Create(&event)

	if err != nil {
		return err
	}

	return nil
}

func ProcessTeamWin(teamToken string) error {
	//Get team token
	token := Token{}
	err := DB.Where("id = ?", teamToken).Last(&token)

	if err != nil {
		return err
	}

	//Get team
	team := Team{}
	err = DB.Where("id = ?", token.ObjectID).Last(&team)

	if err != nil {
		return err
	}

	//Get meeting
	meeting := Meeting{}
	err = DB.Where("f_team_id = ? OR s_team_id = ?", team.ID, team.ID).Last(&meeting)

	if err != nil {
		return err
	}

	//set meeting end_dt
	meeting.EndDt = nulls.Time{
		Time:  time.Now(),
		Valid: true,
	}

	err = DB.Update(&meeting)

	if err != nil {
		return err
	}

	//get tournament token
	tournamentToken := Token{}
	err = DB.Where("object_id = ?", meeting.TournamentID).Last(&tournamentToken)

	if err != nil {
		return err
	}

	//create event with tournament token reciever
	event := Event{
		SenderID:   token.ID,
		ReceiverID: tournamentToken.ID,
		Type:       TEAM_WINS,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = DB.Create(&event)

	if err != nil {
		return err
	}

	return nil
}

func EndMeeting(meetingID int) error {
	//Get meeting
	meeting := Meeting{}
	err := DB.Where("id = ?", meetingID).Last(&meeting)

	if err != nil {
		return err
	}

	//set meeting end_dt
	meeting.EndDt = nulls.Time{
		Time:  time.Now(),
		Valid: true,
	}

	err = DB.Update(&meeting)

	if err != nil {
		return err
	}

	return nil
}
