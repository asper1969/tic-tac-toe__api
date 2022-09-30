package actions

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"tic-tac-toe__api/models"

	"github.com/gobuffalo/buffalo"
)

// QuizCategories default implementation.
func QuizCategories(c buffalo.Context) error {
	categories := models.Categories{}
	err := models.DB.Where("published = true").All(&categories)

	if err != nil {
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(categories))
}

// QuizQuestions default implementation.
func QuizQuestions(c buffalo.Context) error {
	params := c.Params().(url.Values)
	categories := []int{}
	levels := []int{}

	for _, i := range params["category[]"] {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		categories = append(categories, j)
	}

	for _, i := range params["difficulty[]"] {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		levels = append(levels, j)
	}

	questions, err := GetQuestionSet(categories, levels)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(questions))
}

func GetQuestionSet(categories []int, levels []int) (models.Questions, error) {
	questions := models.Questions{}
	dbQuery := models.DB.Where("published = true")

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

	err := dbQuery.EagerPreload().Order("RANDOM()").Limit(100).All(&questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}
