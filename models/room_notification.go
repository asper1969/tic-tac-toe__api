package models

import (
	"encoding/json"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// RoomNotification is used by pop to map your room_notifications database table to your go code.
type RoomNotification struct {
	ID       int    `json:"id" db:"id"`
	Room     string `json:"room" db:"room"`
	Status   int    `json:"status" db:"status"`
	UpdateDt string `json:"update_dt" db:"update_dt"`
	Type     int    `json:"type" db:"type"`
	Data     string `json:"data" db:"data"`
}

// String is not required by pop and may be deleted
func (r RoomNotification) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// RoomNotifications is not required by pop and may be deleted
type RoomNotifications []RoomNotification

// String is not required by pop and may be deleted
func (r RoomNotifications) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *RoomNotification) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *RoomNotification) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *RoomNotification) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
