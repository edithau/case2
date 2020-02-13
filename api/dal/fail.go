package dal

import (
	"errors"

	"github.com/InVisionApp/case-study/api/team"
)

// failDAL is a DAL that only has database failures.
type failDAL struct{}

func NewFailDAL() *failDAL                                { return &failDAL{} }
func (d *failDAL) Authenticate(email, pass string) error  { return errors.New("database error") }
func (d *failDAL) SessionTeams(sessionID string) []string { return nil }
func (d *failDAL) GetTeamByID(teamID string) *team.Team   { return nil }
