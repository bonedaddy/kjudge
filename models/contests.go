package models

import (
	"git.nkagami.me/natsukagami/kjudge/db"
	"git.nkagami.me/natsukagami/kjudge/models/verify"
	"github.com/pkg/errors"
)

// ContestType is the enum representing the contest type.
// The contest type determines ONLY how the scoreboard is rendered.
type ContestType string

const (
	// The contestants are sorted by number of "solved" problems, then by total penalty.
	ContestTypeUnweighted ContestType = "unweighted"
	// Problems each have different scores, and penalty only serves as tiebreakers.
	ContestTypeWeighted ContestType = "weighted"
)

// Verify tries to verify the values of a Contest struct.
func (c *Contest) Verify() error {
	if err := verify.Names(c.Name); err != nil {
		return errors.WithMessage(err, "name: ")
	}
	if !c.StartTime.Before(c.EndTime) {
		return errors.New("start time: must be before end time")
	}
	if c.ContestType != ContestTypeUnweighted &&
		c.ContestType != ContestTypeWeighted {
		return errors.New("contest type: invalid value")
	}
	return nil
}

// GetContestsUnfinished gets a list of contests that are unfinished (upcoming or pending).
func GetContestsUnfinished(db db.DBContext) ([]*Contest, error) {
	var res []*Contest
	if err := db.Select(&res, "SELECT * FROM contests WHERE datetime(end_time) > datetime('now')"+queryContestOrderBy); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

// GetContestsFinished gets a list of contests that are finished (past contests).
func GetContestsFinished(db db.DBContext) ([]*Contest, error) {
	var res []*Contest
	if err := db.Select(&res, "SELECT * FROM contests WHERE datetime(end_time) <= datetime('now')"+queryContestOrderBy); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

// GetContests returns a list of all contests.
func GetContests(db db.DBContext) ([]*Contest, error) {
	var res []*Contest
	if err := db.Select(&res, "SELECT * FROM contests ORDER BY id DESC"); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

// GetNearestOngoingContest return the nearest ongoing contest
func GetNearestOngoingContest(db db.DBContext) (*Contest, error) {
	var res Contest
	if err := db.Get(&res, "SELECT * FROM contests WHERE datetime(end_time) > datetime('now') AND date(start_time) <= datetime('now')"+queryContestOrderBy+" LIMIT 1"); err != nil {
		return nil, errors.WithStack(err)
	}
	return &res, nil
}
