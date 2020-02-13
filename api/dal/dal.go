package dal

import "github.com/InVisionApp/case-study/api/team"

// DAL is a simple interface to manage datastores.
type DAL interface {
	// Authenticate verifies login credentials. It returns nil on success, otherwise an error.
	Authenticate(email string, password string) error
	// SessionTeams returns the valid teams for a session ID, or nil.
	SessionTeams(sessionID string) []string
	// GetTeamByID returns a Team matchin teamID, or nil.
	GetTeamByID(teamID string) *team.Team
}
