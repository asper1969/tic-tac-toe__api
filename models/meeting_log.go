package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// MeetingLog is used by pop to map your meeting_logs database table to your go code.
type MeetingLog struct {
	ID           uuid.UUID `json:"id" db:"id"`
	MeetingID    int       `json:"meeting_id" db:"meeting_id"`
	PlacesSet    string    `json:"places_set" db:"places_set"`
	QuestionsLog string    `json:"quesions_log" db:"questions_log"`
	FTeamScore   int       `json:"f_team_score" db:"f_team_score"`
	STeamScore   int       `json:"s_team_score" db:"s_team_score"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Meeting      *Meeting  `json:"-" belongs_to:"meeting"`
}

// String is not required by pop and may be deleted
func (m MeetingLog) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MeetingLogs is not required by pop and may be deleted
type MeetingLogs []MeetingLog

// String is not required by pop and may be deleted
func (m MeetingLogs) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MeetingLog) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MeetingLog) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MeetingLog) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
