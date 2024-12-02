package request

type LoginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
