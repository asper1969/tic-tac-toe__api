package actions

import (
	"tic-tac-toe__api/locales"
	"tic-tac-toe__api/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n/v2"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app *buffalo.App
	T   *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_tic_tac_toe_api_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.

		app.Use(popmw.Transaction(models.DB))
		app.GET("/", HomeHandler)
		app.GET("/api/categories", QuizCategories)
		app.GET("/api/questions", QuizQuestions)
		app.POST("/api/session", SessionCreate)
		app.GET("/api/session", SessionGet)
		app.PUT("/api/session", SessionUpdate)

		app.GET("/api/admin/categories", QuizCategories)
		app.GET("/api/admin/questions", QuizQuestionsAdmin)
		app.GET("/api/admin/questions/{id}", GetQuestionAdmin)
		//TODO: add admin routes

		// */api/admin/questions/:id - *GET*/POST/PUT/DELETE by :id

		app.POST("/api/tournaments/create", TournamentsCreate)
		app.POST("/api/tournaments/join", TournamentsJoin)
		app.POST("/api/tournaments/create_round", TournamentsCreateRound)
		app.POST("/api/tournaments/start_round", TournamentsStartRound)
		app.POST("/api/tournaments/stop_round", TournamentsStopRound)
		app.POST("/api/tournaments/stop", TournamentsStop)
		app.POST("/api/tournaments/end", TournamentsEnd)
		app.POST("/api/tournaments/pause", TournamentsPause)
		app.POST("/api/tournaments/continue", TournamentsContinue)
		app.POST("/api/tournaments/action", TournamentsAction)
		app.GET("/api/events/get", EventsGet)
		app.POST("/api/meetings/make_move", MeetingsMakeMove)
		app.POST("/api/meetings/answer_question", MeetingsAnswerQuestion)
		app.POST("/api/meetings/accept_move", MeetingsAcceptMove)
		app.POST("/api/meetings/decline_move", MeetingsDeclineMove)
		app.POST("/api/meetings/pass_move", MeetingsPassMove)
		app.POST("/api/meetings/team_win", MeetingsTeamWin)
		app.POST("/api/tournaments/check_token", TournamentTokenIsActive)
		app.GET("/api/tournaments/get_session", TournamentGetSession)

		app.GET("/api/tests/get", TestsGet)
		app.POST("/api/tests/post", TestsPost)
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
