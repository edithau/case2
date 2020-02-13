package session

const (
	FakeSessionID = "56906974-f924-4a1c-889e-a1dd2f395ac2"
)

type Session struct {
	SessionID string   `json:"session_id"`
	TeamIDs   []string `json:"teams"`
}

var fakeSessions = []*Session{
	{
		SessionID: FakeSessionID,
		TeamIDs:   []string{"ch72gsb320000000000000000", "ch72gsb320000000000000001"},
	},
}

func GetAll() []*Session { return fakeSessions }
