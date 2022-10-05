package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tic-tac-toe__api/models"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/thanhpk/randstr"
)

type CreateRequest struct {
	FTeam      string `json:"f_team"`
	STeam      string `json:"s_team"`
	Categories []int  `json:"categories"`
	Levels     []int  `json:"levels"`
	MaxScore   int    `json:"max_score"`
}

type UpdateRequest struct {
	GamePass     string `json:"game_pass"`
	PlacesSet    []int  `json:"places_set"`
	QuestionsLog []int  `json:"questions_log"`
	FTeamScore   int    `json:"f_team_score"`
	STeamScore   int    `json:"s_team_score"`
	TeamID       int    `json:"team_id"`
}

// SessionCreate default implementation.
func SessionCreate(c buffalo.Context) error {
	requestData := &CreateRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	gamePass := generateGamePass()
	questions, err := GetQuestionSet(requestData.Categories, requestData.Levels)

	if err != nil {
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	questionsSet, _ := json.Marshal(questions)
	categories, _ := json.Marshal(requestData.Categories)
	levels, _ := json.Marshal(requestData.Levels)

	session := models.Session{
		FTeam:        requestData.FTeam,
		STeam:        requestData.STeam,
		Categories:   string(categories),
		Levels:       string(levels),
		GamePass:     gamePass,
		StartDt:      time.Now().Format("2006-02-01 00:00:00"),
		QuestionsSet: string(questionsSet),
		MaxScore:     requestData.MaxScore,
	}

	err = models.DB.Create(&session)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(session))
}

func generateGamePass() string {
	str := randstr.String(5)
	session := models.Session{}
	err := models.DB.Where("game_pass = ?", str).First(&session)

	if err == nil {
		str = generateGamePass()
	}

	return str
}

// SessionGet default implementation.
func SessionGet(c buffalo.Context) error {
	gamePass := c.Param("game_pass")
	session := models.Session{}

	err := models.DB.Where("game_pass = ?", gamePass).First(&session)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(session))
}

// SessionUpdate default implementation.
func SessionUpdate(c buffalo.Context) error {
	requestData := &UpdateRequest{}

	if err := c.Bind(requestData); err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Get session by game_pass
	session := models.Session{}
	err := models.DB.Where("game_pass = ?", requestData.GamePass).First(&session)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	//Add record with actual data to sessions_log
	placesSet, _ := json.Marshal(requestData.PlacesSet)
	questionsLog, _ := json.Marshal(requestData.QuestionsLog)

	sessionLog := models.SessionLog{
		UpdateDt:     time.Now().Format("2006-02-01 00:00:00"),
		PlacesSet:    string(placesSet),
		QuestionsLog: string(questionsLog),
		FTeamScore:   requestData.FTeamScore,
		STeamScore:   requestData.STeamScore,
		Session:      &session,
	}

	err = models.DB.Create(&sessionLog)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	//Add record to room_notifications
	roomNotification := models.RoomNotification{
		Room:     fmt.Sprintf("%s:%d", requestData.GamePass, requestData.TeamID),
		Status:   1,
		UpdateDt: time.Now().Format("2006-02-01 00:00:00"),
	}

	err = models.DB.Create(&roomNotification)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusBadGateway, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(session))
}
