package domain

type Chat struct {
	ID             string   `json:"ID"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	NameMessagesDB string   `json:"nameMessegesDB"`
}
