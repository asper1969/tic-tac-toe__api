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
	ID           int        `json:"id" db:"id"`
	GamePass     string     `json:"game_pass" db:"game_pass"`
	Locale       string     `json:"locale" db:"locale"`
	StartDt      nulls.Time `json:"start_dt" db:"start_dt"`
	EndDt        nulls.Time `json:"end_dt" db:"end_dt"`
	MaxScore     int        `json:"max_score" db:"max_score"`
	TeamsAmount  int        `json:"teams_amount" db:"teams_amount"`
	FieldsAmount int        `json:"fields_amount" db:"fields_amount"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type MeetingPair struct {
	Team  *Team
	Field int
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

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("!!!!!!!!!!!!!!panic occurred:", err)
	// 	}
	// }()

	//Get last round
	//Create meetings for next round
	//--Get all teams and create meetings for each pair
	lastMeeting := Meeting{}
	err := DB.Where("tournament_id = ?", t.ID).Last(&lastMeeting)
	round := 1

	//No meetings in tournament
	if err == nil {
		round = lastMeeting.Round + 1
	}

	//Create new meetings for first round
	//Get all tournament teams
	teams := Teams{}
	err = DB.Where("tournament_id = ?", t.ID).Order("id asc").All(&teams)

	if err != nil {
		fmt.Println(err)
		return err
	}

	opponents := map[int][]Team{}

	//*********
	for i, team := range teams {
		teamField := GetMeetingField(i+1, len(teams), t.FieldsAmount, round)

		if teamField > 0 {
			opponents[teamField] = append(opponents[teamField], team)
		}
	}

	for field, pair := range opponents {
		firstTeam := pair[0]
		secondTeam := pair[1]
		meeting := Meeting{
			FTeamID:      firstTeam.ID,
			STeamID:      secondTeam.ID,
			TournamentID: t.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Round:        round,
			Field:        field,
		}

		//TODO: add rounds logic
		if round == 1 {
			questionsSet, err := GetQuestionSet([]int{}, []int{}, t.Locale)

			if err != nil {
				fmt.Println(err)
				return err
			}

			questionsSetStr, _ := json.Marshal(questionsSet)
			meeting.QuestionsSet = string(questionsSetStr)
		} else {
			meeting.QuestionsSet = "[]"
		}

		err = DB.Create(&meeting)

		if err != nil {
			return err
		}
	}

	//*********

	// mid := len(teams) / 2
	// t1 := teams[:mid]
	// t2 := teams[mid:]

	// for i := 0; i < len(t1); i++ {
	// 	firstTeam := t1[i]
	// 	secondTeam := t2[i]

	// 	//TODO: calculate by previous meetings
	// 	//For round get last meeting round and add to it 1, if there are not another meetings, then round = 1
	// 	round := 1
	// 	//For field get last meeting field and calculate field number by alghorithm
	// 	//number of teams, number of fields and round for calculating
	// 	field := 1

	// 	meeting := Meeting{
	// 		FTeamID:      firstTeam.ID,
	// 		STeamID:      secondTeam.ID,
	// 		TournamentID: t.ID,
	// 		CreatedAt:    time.Now(),
	// 		UpdatedAt:    time.Now(),
	// 		Round:        round,
	// 		Field:        field,
	// 	}

	// 	//TODO: add rounds logic
	// 	if round == 1 {
	// 		questionsSet, err := GetQuestionSet([]int{}, []int{}, t.Locale)

	// 		if err != nil {
	// 			return err
	// 		}

	// 		questionsSetStr, _ := json.Marshal(questionsSet)
	// 		meeting.QuestionsSet = string(questionsSetStr)
	// 	}

	// 	err = DB.Create(&meeting)

	// 	if err != nil {
	// 		return err
	// 	}
	// }

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
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

var TournamentTable = map[int]interface{}{
	2: map[int]interface{}{ //2 teams
		1: map[int]interface{}{ //Court 1
			1: [][]int{{1, 2}},
			2: [][]int{{1, 2}},
			3: [][]int{{1, 2}},
			4: [][]int{{1, 2}},
			5: [][]int{{1, 2}},
			6: [][]int{{1, 2}},
		},
		2: map[int]interface{}{ //Court 2
			1: [][]int{{1, 2}, {}}, //Round 1: Court 1: t1-t2, Court 2: Empty
			2: [][]int{{}, {1, 2}},
			3: [][]int{{1, 2}, {}},
			4: [][]int{{}, {1, 2}},
			5: [][]int{{1, 2}, {}},
			6: [][]int{{}, {1, 2}},
		},
		3: map[int]interface{}{ //Court 3
			1: [][]int{{1, 2}, {}, {}}, //Round 1: Court 1: t1-t2, Court 2: Empty, Court 3: Empty
			2: [][]int{{}, {1, 2}, {}},
			3: [][]int{{}, {}, {1, 2}},
			4: [][]int{{1, 2}, {}, {}},
			5: [][]int{{}, {1, 2}, {}},
			6: [][]int{{}, {}, {1, 2}},
		},
	},
	3: map[int]interface{}{ //3 teams
		1: map[int]interface{}{
			1: [][]int{{1, 2}},
			2: [][]int{{2, 3}},
			3: [][]int{{1, 3}},
			4: [][]int{{1, 2}},
			5: [][]int{{2, 3}},
			6: [][]int{{1, 3}},
		},
		2: map[int]interface{}{
			1: [][]int{{1, 2}, {}}, //Round 1: Court 1: t1-t2, Court 2: Empty
			2: [][]int{{}, {2, 3}},
			3: [][]int{{1, 3}, {}},
			4: [][]int{{}, {1, 2}},
			5: [][]int{{2, 3}, {}},
			6: [][]int{{}, {1, 3}},
		},
		3: map[int]interface{}{
			1: [][]int{{1, 2}, {}, {}}, //Round 1: Court 1: t1-t2, Court 2: Empty, Court 3: Empty
			2: [][]int{{}, {2, 3}, {}},
			3: [][]int{{}, {}, {1, 3}},
			4: [][]int{{2, 3}, {}, {}},
			5: [][]int{{}, {1, 3}, {}},
			6: [][]int{{}, {}, {1, 2}},
		},
	},
	4: map[int]interface{}{ //4 teams
		1: map[int]interface{}{
			1: [][]int{{1, 2}},
			2: [][]int{{3, 4}},
			3: [][]int{{1, 3}},
			4: [][]int{{2, 4}},
			5: [][]int{{1, 4}},
			6: [][]int{{2, 3}},
		},
		2: map[int]interface{}{
			1: [][]int{{1, 2}, {3, 4}}, //Round 1: Court 1: t1-t2, Court 2: t3-t4
			2: [][]int{{1, 3}, {2, 4}},
			3: [][]int{{1, 4}, {2, 3}},
			4: [][]int{{2, 3}, {1, 4}},
			5: [][]int{{3, 4}, {1, 2}},
			6: [][]int{{2, 4}, {1, 3}},
		},
		3: map[int]interface{}{
			1: [][]int{{1, 2}, {3, 4}, {}}, //Round 1: Court 1: t1-t2, Court 2: t3-t4, Court 3: Empty
			2: [][]int{{}, {1, 3}, {2, 4}},
			3: [][]int{{1, 4}, {}, {2, 3}},
			4: [][]int{{2, 3}, {1, 4}, {}},
			5: [][]int{{}, {2, 4}, {1, 3}},
			6: [][]int{{2, 4}, {}, {1, 3}},
		},
	},
	6: map[int]interface{}{ //6 teams
		1: map[int]interface{}{
			1: [][]int{{1, 2}},
			2: [][]int{{3, 4}},
			3: [][]int{{5, 6}},
			4: [][]int{{1, 3}},
			5: [][]int{{2, 5}},
			6: [][]int{{4, 6}},
		},
		2: map[int]interface{}{
			1: [][]int{{1, 2}, {3, 4}}, //Round 1: Court 1: t1-t2, Court 2: t3-t4
			2: [][]int{{1, 5}, {2, 6}},
			3: [][]int{{2, 3}, {4, 5}},
			4: [][]int{{1, 4}, {5, 6}},
			5: [][]int{{3, 6}, {1, 4}},
			6: [][]int{{2, 5}, {3, 6}},
		},
		3: map[int]interface{}{
			1: [][]int{{1, 2}, {3, 4}, {5, 6}}, //Round 1: Court 1: t1-t2, Court 2: t3-t4, Court 3: t5-t6
			2: [][]int{{3, 5}, {2, 6}, {1, 4}},
			3: [][]int{{3, 6}, {5, 1}, {2, 4}},
			4: [][]int{{1, 6}, {2, 4}, {3, 5}},
			5: [][]int{{4, 5}, {3, 6}, {1, 2}},
			6: [][]int{{1, 3}, {2, 5}, {4, 6}},
		},
	},
}

func GetMeetingField(teamIdx int, teamsAmount int, fieldsAmount int, currentRound int) int {
	opponents, ok := TournamentTable[teamsAmount].(map[int]interface{})[fieldsAmount].(map[int]interface{})[currentRound].([][]int)

	if !ok {
		return 0
	}

	for i, pair := range opponents {

		if len(pair) > 0 && (pair[0] == teamIdx || pair[1] == teamIdx) {
			return i + 1
		}
	}

	return 0
}
