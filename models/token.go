package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

type TokenType int

const (
	TOKEN_TOURNAMENT TokenType = 1
	TOKEN_MODERATOR  TokenType = 2
	TOKEN_TEAM       TokenType = 3
	TOKEN_FIELD      TokenType = 4
)

// Token is used by pop to map your tokens database table to your go code.
type Token struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Type      TokenType `json:"type" db:"type"`
	ObjectID  int       `json:"object_id" db:"object_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Expired   bool      `json:"expired" db:"expired"`
}

// String is not required by pop and may be deleted
func (t Token) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tokens is not required by pop and may be deleted
type Tokens []Token

// String is not required by pop and may be deleted
func (t Tokens) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Token) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Token) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Token) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
