package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// Tournament is used by pop to map your tournaments database table to your go code.
type Tournament struct {
	ID        int        `json:"id" db:"id"`
	GamePass  string     `json:"game_pass" db:"game_pass"`
	Locale    string     `json:"locale" db:"locale"`
	StartDt   nulls.Time `json:"start_dt" db:"start_dt"`
	EndDt     nulls.Time `json:"end_dt" db:"end_dt"`
	MaxScore  int        `json:"max_score" db:"max_score"`
	Rounds    int        `json:"rounds" db:"rounds"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (t Tournament) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tournaments is not required by pop and may be deleted
type Tournaments []Tournament

// String is not required by pop and may be deleted
func (t Tournaments) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tournament) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Tournament) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Tournament) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (t *Tournament) GetToken() Token {
	token := Token{}
	err := DB.Where("object_id = ?", t.ID).First(&token)

	if err != nil {
		fmt.Println("Token for object not found")
	}

	return token
}
