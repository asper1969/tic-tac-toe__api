package models

import (
	"encoding/json"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// Session is used by pop to map your sessions database table to your go code.
type Session struct {
	ID           int          `json:"id" db:"id"`
	FTeam        string       `json:"f_team" db:"f_team"`
	STeam        string       `json:"s_team" db:"s_team"`
	GamePass     string       `json:"game_pass" db:"game_pass"`
	QuestionsSet string       `json:"questions_set" db:"questions_set"`
	Levels       string       `json:"levels" db:"levels"`
	Categories   string       `json:"categories" db:"categories"`
	MaxScore     int          `json:"max_score" db:"max_score"`
	StartDt      string       `json:"start_dt" db:"start_dt"`
	EndDt        nulls.String `json:"-" db:"end_dt"`
	SessionLogs  []SessionLog `json:"session_logs" has_many:"session_logs"`
}

// String is not required by pop and may be deleted
func (s Session) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Sessions is not required by pop and may be deleted
type Sessions []Session

// String is not required by pop and may be deleted
func (s Sessions) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Session) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Session) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Session) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
