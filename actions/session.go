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
	// _, err := GetQuestionSet(requestData.Categories, requestData.Levels)

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
		StartDt:      time.Now(),
		QuestionsSet: string(questionsSet),
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
