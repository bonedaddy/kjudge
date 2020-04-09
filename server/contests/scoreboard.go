package contests

import (
	"net/http"
	"time"

	"git.nkagami.me/natsukagami/kjudge/db"
	"git.nkagami.me/natsukagami/kjudge/models"
	"git.nkagami.me/natsukagami/kjudge/server/user"
	"github.com/labstack/echo/v4"
)

// ScoreboardCtx is the context required to display the scoreboard page
type ScoreboardCtx struct {
	*user.AuthCtx
	*models.Scoreboard
}

// Render renders the scoreboard context
func (s *ScoreboardCtx) Render(c echo.Context) error {
	return c.Render(http.StatusOK, "contests/scoreboard", s)
}

// RenderJSON renders a scoreboard in JSON.
func (s *ScoreboardCtx) RenderJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, s.JSON())
}

// Collect a ScoreboardCtx
func getScoreboardCtx(db db.DBContext, c echo.Context) (*ScoreboardCtx, error) {
	contestCtx, err := getContestCtx(db, c)
	if err != nil {
		return nil, err
	}

	// get contest's problems
	contest := contestCtx.Contest
	// get contest information
	problems := contestCtx.Problems

	scoreboard, err := models.GetScoreboard(db, contest, problems)
	if err != nil {
		return nil, err
	}

	return &ScoreboardCtx{
		AuthCtx:    contestCtx.AuthCtx,
		Scoreboard: scoreboard,
	}, nil
}

// ScoreboardGet implements GET /contest/:id/scoreboard
func (g *Group) ScoreboardGet(c echo.Context) error {
	ctx, err := getScoreboardCtx(g.db, c)
	if err != nil {
		return err
	}
	return ctx.Render(c)
}

// ScoreboardJSONGet implements GET /contest/:id/scoreboard/json
func (g *Group) ScoreboardJSONGet(c echo.Context) error {
	ctx, err := getScoreboardCtx(g.db, c)
	if err != nil {
		return err
	}
	// if contest has ended, scoreboard should be accessible to everyone
	if ctx.Contest.EndTime.Before(time.Now()) {
		return ctx.RenderJSON(c)
	}

	// allow users to access scoreboard JSON based on ScoreboardViewStatus
	if ctx.Contest.ScoreboardViewStatus == models.ScoreboardViewStatusNoScoreboard ||
		ctx.Contest.ScoreboardViewStatus == models.ScoreboardViewStatusUser && ctx.AuthCtx.Me == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Scoreboard access is not granted")
	} else {
		return ctx.RenderJSON(c)
	}
}
