package dal

import (
	"github.com/InVisionApp/interview-test/api/session"
	"github.com/InVisionApp/interview-test/api/team"
	"github.com/InVisionApp/interview-test/api/user"
	"golang.org/x/crypto/bcrypt"
)

// mockDAL pretends to be a working DAL.
type mockDAL struct{}

func NewMockDAL() *mockDAL { return &mockDAL{} }

func (d *mockDAL) Authenticate(email, pass string) error {
	users := user.GetAll()
	for i := range users {
		if users[i].Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(users[i].Hash), []byte(pass))
			if err != nil {
				return &UnauthorizedErr{}
			}
			return nil
		}
	}
	return &CredentialsErr{}
}

func (d *mockDAL) SessionTeams(sessionID string) []string {
	sessions := session.GetAll()
	for i := range sessions {
		if sessions[i].SessionID == sessionID {
			return sessions[i].TeamIDs
		}
	}
	return nil
}

func (d *mockDAL) GetTeamByID(teamID string) *team.Team {
	teams := team.GetAll()
	for i := range teams {
		if teams[i].TeamID == teamID {
			return teams[i]
		}
	}
	return nil
}
