package actions

import (
	"fmt"
	"net/http"
	"net/url"
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
	questions := models.Questions{}
	dbQuery := models.DB.Where("published = true")

	params := c.Params().(url.Values)
	categories := params["category[]"]
	levels := params["difficulty[]"]

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

	err := dbQuery.EagerPreload().All(&questions)

	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, r.JSON(err))
	}

	return c.Render(http.StatusOK, r.JSON(questions))
}
