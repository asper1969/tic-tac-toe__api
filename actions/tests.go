package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// TestsPost default implementation.
func TestsPost(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("tests/post.html"))
}

