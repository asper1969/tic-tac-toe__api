package models

import (
	"encoding/json"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// Answer is used by pop to map your answers database table to your go code.
type Answer struct {
	ID         int       `json:"id" db:"id"`
	Text       string    `json:"text" db:"text"`
	QuestionID int       `json:"-" db:"question_id"`
	IsRight    bool      `json:"is_right" db:"is_right"`
	Hash       string    `json:"hash" db:"hash"`
	Question   *Question `json:"-" belongs_to:"question"`
}

// String is not required by pop and may be deleted
func (a Answer) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Answers is not required by pop and may be deleted
type Answers []Answer

// String is not required by pop and may be deleted
func (a Answers) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Answer) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Answer) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Answer) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
