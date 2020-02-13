package user

type User struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Hash   string `json:"-"` // Password bcrypt string hash.
}

// List of fake users that have access.
var fakeUsers = []*User{
	{
		UserID: 1,
		Email:  "user@example.com",
		Hash:   "$2a$10$CI3AFyP8DQq/EXH3.skPdOrpHdOMud.U1g6XDXgdREqXOCImvzArS",
	},
	{
		UserID: 1337,
		Email:  "baduser@example.com",
		Hash:   "$2a$10$WPw.6PA6LAqy4UM1BK8f1.tzYCN0rL8erKQtKVZyrZN9k6BoA9Zdu",
	},
}

func GetAll() []*User { return fakeUsers }
