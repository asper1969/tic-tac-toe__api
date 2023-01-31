package models

import (
	"encoding/json"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// Question is used by pop to map your questions database table to your go code.
type Question struct {
	ID         int          `json:"id" db:"id"`
	Difficulty int          `json:"difficulty" db:"difficulty"`
	Text       string       `json:"text" db:"text"`
	CategoryID int          `json:"category_id" db:"category_id"`
	Published  bool         `json:"published" db:"published"`
	Hash       nulls.String `json:"hash" db:"hash"`
	Locale     string       `json:"locale" db:"locale"`
	Answers    []Answer     `json:"answers" has_many:"answers"`
}

// String is not required by pop and may be deleted
func (q Question) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Questions is not required by pop and may be deleted
type Questions []Question

// String is not required by pop and may be deleted
func (q Questions) String() string {
	jq, _ := json.Marshal(q)
	return string(jq)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (q *Question) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (q *Question) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (q *Question) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func GetQuestionSet(categories []int, levels []int, locale string) (Questions, error) {
	questions := Questions{}
	dbQuery := DB.Where("published = true").Where("locale = ?", locale)

	if len(categories) > 0 {
		/**
		* Get all published filtered by categories
		**/
		dbQuery = dbQuery.Where("category_id IN (?)", categories)
	}

	if len(levels) > 0 {
		/**
		* Get all published filtered by levels
		**/
		dbQuery = dbQuery.Where("difficulty IN (?)", levels)
	}

	err := dbQuery.EagerPreload().Order("RAND()").Limit(100).All(&questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}
