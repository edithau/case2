package team

type Team struct {
	TeamID  string `json:"team_id"`
	Name    string `json:"name"`
	Members int    `json:"members"`
}

var fakeTeams = []*Team{
	{
		TeamID:  "ch72gsb320000000000000000",
		Name:    "Team A",
		Members: 203,
	},
	{
		TeamID:  "ch72gsb320000000000000001",
		Name:    "Team B",
		Members: 5,
	},
}

func GetAll() []*Team { return fakeTeams }
