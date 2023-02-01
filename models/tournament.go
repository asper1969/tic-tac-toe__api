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

func (t *Tournament) CreateNextRound() error {
	//Get last round
	//Create meetings for next round
	//--Get all teams and create meetings for each pair
	lastMeeting := Meeting{}
	err := DB.Where("tournament_id = ?", t.ID).Last(&lastMeeting)

	//No meetings in tournament
	if err != nil {
		//Create new meetings for first round
		//Get all tournament teams
		teams := Teams{}
		err = DB.Where("tournament_id = ?", t.ID).All(&teams)

		if err != nil {
			return err
		}

		mid := len(teams) / 2
		t1 := teams[:mid]
		t2 := teams[mid:]

		for i := 0; i < len(t1); i++ {
			firstTeam := t1[i]
			secondTeam := t2[i]

			round := 1

			meeting := Meeting{
				FTeamID:      firstTeam.ID,
				STeamID:      secondTeam.ID,
				TournamentID: t.ID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Round:        round,
			}

			//TODO: add rounds logic
			if round == 1 {
				questionsSet, err := GetQuestionSet([]int{}, []int{}, t.Locale)

				if err != nil {
					return err
				}

				questionsSetStr, _ := json.Marshal(questionsSet)
				meeting.QuestionsSet = string(questionsSetStr)
			}

			err = DB.Create(&meeting)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Tournament) StartNextRound() error {
	//Get all new tournament meetings (start_dt == NULL)
	meetings := Meetings{}
	err := DB.Where("start_dt IS NULL AND tournament_id = ?", t.ID).All(&meetings)

	if err != nil {
		return err
	}

	//Set start_dt for each meeting
	for _, meeting := range meetings {
		meeting.StartDt = nulls.Time{
			Time:  time.Now(),
			Valid: true,
		}
		DB.Update(&meeting)

		teams := Teams{}
		err := DB.Where("id = ? || id = ?", meeting.FTeamID, meeting.STeamID).All(&teams)

		if err != nil {
			return err
		}

		//Create ROUND_START events for each meeting teams
		for _, team := range teams {
			//Get team token
			token := Token{}
			err := DB.Where("object_id = ?", team.ID).Last(&token)
			fmt.Println(team.ID)

			if err != nil {
				return err
			}

			event := Event{
				SenderID:   t.GetToken().ID,
				ReceiverID: token.ID,
				Type:       ROUND_START,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			err = DB.Create(&event)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
