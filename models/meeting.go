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
	Teams        []Team       `json:"teams" has_many:"teams"`
	MeetingLogs  []MeetingLog `json:"meeting_logs" has_many:"meeting_logs"`
	Round        int          `json:"round" db:"round"`
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
