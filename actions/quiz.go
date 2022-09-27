package actions

import (
	"net/http"
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
