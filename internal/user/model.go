package user

type User struct {
	ID       string `json:"ID"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
