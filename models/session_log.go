package models

import (
	"encoding/json"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// SessionLog is used by pop to map your session_logs database table to your go code.
type SessionLog struct {
	ID           int      `json:"id" db:"id"`
	SessionID    int      `json:"session_id" db:"session_id"`
	UpdateDt     string   `json:"update_dt" db:"update_dt"`
	PlacesSet    string   `json:"places_set" db:"places_set"`
	QuestionsLog string   `json:"questions_log" db:"questions_log"`
	FTeamScore   int      `json:"f_team_score" db:"f_team_score"`
	STeamScore   int      `json:"s_team_score" db:"s_team_score"`
	Session      *Session `json:"-" belongs_to:"session"`
}

// String is not required by pop and may be deleted
func (s SessionLog) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// SessionLogs is not required by pop and may be deleted
type SessionLogs []SessionLog

// String is not required by pop and may be deleted
func (s SessionLogs) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *SessionLog) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *SessionLog) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *SessionLog) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
