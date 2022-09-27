package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// QuizCategories default implementation.
func QuizCategories(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Categories list"}))
}
